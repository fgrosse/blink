package blink

import (
	"io"
	"testing"
)

func TestLEDImplementsIoCloser(t *testing.T) {
	var _ io.Closer = new(LED)
}

func TestDontPanicWhenClosingNilLEDs(t *testing.T) {
	var led *LED
	led.Close()
}
