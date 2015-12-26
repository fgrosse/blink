package blink

import (
	"errors"
	"fmt"
)

const (
	VendorNumber  = 0x27B8
	ProductNumber = 0x01ED
)

var NoDeviceErr = errors.New("could not find blink1 device")

type LED struct {
	*usbDevice
	info usbDeviceInfo
}

func New() (*LED, error) {
	l := LED{}

	found := false
	for di := range usbDevices() {
		if di.vendorId == VendorNumber && di.productId == ProductNumber {
			l.info = di
			found = true
			break
		}
	}

	if !found {
		return nil, NoDeviceErr
	}

	var err error
	l.usbDevice, err = l.info.open()
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
	buf := make([]byte, 8)

	buf[0] = 1
	buf[1] = 'n'
	buf[2] = r
	buf[3] = g
	buf[4] = b
	buf[5] = 0
	buf[6] = 0
	buf[7] = 0

	fmt.Printf("Sending: %#x\n", buf)
	return l.write(buf)
}
