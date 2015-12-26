package blink

import (
	"errors"
	"fmt"
	"time"
)

const (
	VendorNumber  = 0x27B8
	ProductNumber = 0x01ED
)

var NoDeviceErr = errors.New("could not find blink1 device")

type LED struct {
	*usbDevice

	ID byte // ID signals which LED to address: 0=all, 1=led#1, 2=led#2, etc. (mk2 only)
}

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

func (l *LED) Close() {
	if l.usbDevice != nil {
		l.close()
	}
}

func (l *LED) SetRGB(r, g, b byte) error {
	return l.write(&setRGBCommand{rgb{r, g, b}})
}

func (l *LED) FadeRGB(r, g, b byte, d time.Duration) error {
	return l.write(&fadeRGBCommand{
		rgb: rgb{r, g, b},
		duration: d,
		n: l.ID,
	})
}
