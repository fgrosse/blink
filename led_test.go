package blink

import (
	"io"
	"testing"
)

func TestLEDImplementsIoCloser(t *testing.T) {
	var _ io.Closer = LED{}
}
