package blink

// #cgo pkg-config: libusb-1.0
// #cgo LDFLAGS: -lusb-1.0
// #include <libusb-1.0/libusb.h>
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

const (
	hidRequestTypeClass   = (0x01 << 5)
	hidRecipientInterface = 0x01
	hidEndpointOut        = 0x00
	hidSetReport          = 0x09
)

func init() {
	C.libusb_init(nil)
}

type usbDeviceInfo struct {
	path                string
	vendorId, productId uint16
}

type usbDevice struct {
	handle *C.libusb_device_handle
	info   usbDeviceInfo
}

// pretty much based on github.com/boombuler/hid
// https://github.com/boombuler/hid/blob/08a7959390cac69dfd373882ac0a4435765a2545/hid_linux.go#L26
func usbDevices() <-chan usbDeviceInfo {
	result := make(chan usbDeviceInfo)
	go func() {
		var devices **C.struct_libusb_device
		count := C.libusb_get_device_list(nil, &devices)
		if count < 0 {
			close(result)
			return
		}
		defer C.libusb_free_device_list(devices, 1)

		for _, dev := range slice(devices, count) {
			di, err := readDeviceInfo(dev)
			if err != nil {
				continue
			}
			result <- di
		}

		close(result)
	}()

	return result
}

// pretty much based on github.com/boombuler/hid
// https://github.com/boombuler/hid/blob/08a7959390cac69dfd373882ac0a4435765a2545/hid_linux.go#L58
func (di usbDeviceInfo) open() (*usbDevice, error) {
	var devices **C.struct_libusb_device
	cnt := C.libusb_get_device_list(nil, &devices)
	if cnt < 0 {
		return nil, fmt.Errorf("could not open usb device with path %q: could not list USB devices", di.path)
	}
	defer C.libusb_free_device_list(devices, 1)

	for _, dev := range slice(devices, cnt) {
		candidate, err := readDeviceInfo(dev)
		if err != nil {
			continue
		}
		if di.path == candidate.path {
			var handle *C.libusb_device_handle
			err := C.libusb_open(dev, &handle)
			if err != 0 {
				return nil, usbError(err)
			}
			dev := &usbDevice{
				info:   candidate,
				handle: handle,
			}
			return dev, nil
		}
	}

	return nil, errors.New("couldn't open USB device")
}

func readDeviceInfo(dev *C.libusb_device) (di usbDeviceInfo, err error) {
	var desc C.struct_libusb_device_descriptor
	if result := C.libusb_get_device_descriptor(dev, &desc); result < 0 {
		return di, usbError(result)
	}

	numbers, err := getPortNumbers(dev)
	if err != nil {
		return di, err
	}

	return usbDeviceInfo{
		path:      fmt.Sprintf("%.4x:%.4x:%s", desc.idVendor, desc.idProduct, numbers),
		vendorId:  uint16(desc.idVendor),
		productId: uint16(desc.idProduct),
	}, nil
}

// entirely based on github.com/boombuler/hid
// https://github.com/boombuler/hid/blob/08a7959390cac69dfd373882ac0a4435765a2545/hid_linux.go#L237
func getPortNumbers(dev *C.libusb_device) (string, error) {
	const maxlen = 7 // As per the USB 3.0 specs, the current maximum limit for the depth is 7
	var numarr [maxlen]C.uint8_t
	len := C.libusb_get_port_numbers(dev, &numarr[0], maxlen)
	if len < 0 || len > maxlen {
		return "", usbError(len)
	}

	var numstr []string = make([]string, len)
	for i := 0; i < int(len); i++ {
		numstr[i] = fmt.Sprintf("%.2x", numarr[i])
	}

	return strings.Join(numstr, "."), nil
}

// entirely based on github.com/boombuler/hid
// https://github.com/boombuler/hid/blob/08a7959390cac69dfd373882ac0a4435765a2545/hid_linux.go#L251
func slice(devices **C.struct_libusb_device, cnt C.ssize_t) []*C.libusb_device {
	var slice []*C.libusb_device
	*(*reflect.SliceHeader)(unsafe.Pointer(&slice)) = reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(devices)),
		Len:  int(cnt),
		Cap:  int(cnt),
	}
	return slice
}

func (d *usbDevice) write(c command) error {
	if d.handle == nil {
		return errors.New("usb device has not been opend")
	}

	const reportId = 1
	const index = 0
	const timeout = 5000

	data := c.bytes()
	n := len(data)
	written := C.libusb_control_transfer(d.handle,
		C.uint8_t(hidEndpointOut|hidRecipientInterface|hidRequestTypeClass),
		C.uint8_t(hidSetReport),
		C.uint16_t(reportId),
		C.uint16_t(index),
		(*C.uchar)(&data[0]),
		C.uint16_t(n),
		C.uint(timeout),
	)

	if int(written) == n {
		return nil
	}

	return usbError(written)
}

func (d *usbDevice) close() {
	if d == nil || d.handle == nil {
		return
	}

	C.libusb_close(d.handle)
	d.handle = nil
}

type usbError C.int

func (e usbError) Error() string {
	return fmt.Sprintf("libusb: %s [code %d]", usbErrorString[e], int(e))
}

var usbErrorString = map[usbError]string{
	C.LIBUSB_SUCCESS:             "success",
	C.LIBUSB_ERROR_IO:            "i/o error",
	C.LIBUSB_ERROR_INVALID_PARAM: "invalid param",
	C.LIBUSB_ERROR_ACCESS:        "bad access",
	C.LIBUSB_ERROR_NO_DEVICE:     "no device",
	C.LIBUSB_ERROR_NOT_FOUND:     "not found",
	C.LIBUSB_ERROR_BUSY:          "device or resource busy",
	C.LIBUSB_ERROR_TIMEOUT:       "timeout",
	C.LIBUSB_ERROR_OVERFLOW:      "overflow",
	C.LIBUSB_ERROR_PIPE:          "pipe error",
	C.LIBUSB_ERROR_INTERRUPTED:   "interrupted",
	C.LIBUSB_ERROR_NO_MEM:        "out of memory",
	C.LIBUSB_ERROR_NOT_SUPPORTED: "not supported",
	C.LIBUSB_ERROR_OTHER:         "unknown error",
}
