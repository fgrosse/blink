package blink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetRGBCommand(t *testing.T) {
	c := setRGBCommand{Color{2, 3, 4}}
	assert.Equal(t, []byte{0x01, 'n', 0x02, 0x03, 0x04, 0x00, 0x00, 0x00}, c.bytes())
}

func TestFadeRGBCommand(t *testing.T) {
	c := fadeRGBCommand{Color: Color{2, 3, 4}, duration: 328010 * time.Millisecond, n: 1}
	assert.Equal(t, []byte{0x01, 'c', 0x02, 0x03, 0x04, 0x80, 0x21, 0x01}, c.bytes())

	c = fadeRGBCommand{Color: Color{255, 255, 255}, duration: 5 * time.Second}
	assert.Equal(t, []byte{0x01, 'c', 0xff, 0xff, 0xff, 0x01, 0xf4, 0x00}, c.bytes())

}

func TestReadRGBCommand(t *testing.T) {
	c := readRGBCommand{n: 1}
	assert.Equal(t, []byte{0x01, 'r', 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, c.bytes())
}
