package blink

import (
	"errors"
	"fmt"
	"time"
)

const (
	// The VendorNumber as defined in github.com/todbot/blink1
	// https://github.com/todbot/blink1/blob/3c51231d302d7676c50f28debaf82adc5bfa9460/commandline/blink1-lib.h#L30
	VendorNumber = 0x27B8

	// The ProductNumber as defined in github.com/todbot/blink1
	// https://github.com/todbot/blink1/blob/3c51231d302d7676c50f28debaf82adc5bfa9460/commandline/blink1-lib.h#L31
	ProductNumber = 0x01ED
)

// NoDeviceErr is the error that New() returns if no connected blink(1) device was found.
var NoDeviceErr = errors.New("could not find blink1 device")

// LED represents a localy connected blink(1) USB device.
type LED struct {
	*usbDevice

	ID byte // ID signals which LED to address: 0=all, 1=led#1, 2=led#2, etc. (mk2 only)
}

// New connects to a locally connected blink(1) USB device.
// The caller must call Close when it is done with this LED.
// The function returns a NoDeviceErr if no blink(1) device can be found
// or another error if the device could not be opened.
//
// Multiple connected devices are not yet supported and the library will just
// pick one of them to talk to.
func New() (*LED, error) {
	l := LED{}

	found := false
	var di usbDeviceInfo
	for di = range usbDevices() {
		if di.vendorId == VendorNumber && di.productId == ProductNumber {
			found = true
			break
		}
	}
	if !found {
		return nil, NoDeviceErr
	}

	var err error
	l.usbDevice, err = di.open()
	if err != nil {
		return nil, fmt.Errorf("could not open blink1 device: %s", err)
	}

	return &l, nil
}

// Close implements io.Closer by closing to the connection to the USB device.
// The function is idempotent and can be called on already closed or never
// opened devices.
func (l *LED) Close() error {
	if l.usbDevice != nil {
		l.close()
	}

	return nil
}

// SetRGB lights up the blink(1) with the specified RGB color immediately.
func (l *LED) SetRGB(r, g, b byte) error {
	return l.write(&setRGBCommand{rgb{r, g, b}})
}

// FadeRGB lights up the blink(1) with the specified RGB color, fading to that color over a specified duration.
func (l *LED) FadeRGB(r, g, b byte, d time.Duration) error {
	return l.write(&fadeRGBCommand{
		rgb:      rgb{r, g, b},
		duration: d,
		n:        l.ID,
	})
}
