package blink_test

import (
	"testing"

	"github.com/fgrosse/blink"
	"github.com/stretchr/testify/assert"
)

func TestParseColor(t *testing.T) {
	var c blink.Color
	c = blink.MustParseColor("123,45,67")
	assert.Equal(t, blink.Color{R: 123, G: 45, B: 67}, c)

	c = blink.MustParseColor("123, 45,  \t67")
	assert.Equal(t, blink.Color{R: 123, G: 45, B: 67}, c)

	c = blink.MustParseColor("#123456")
	assert.Equal(t, blink.Color{R: 18, G: 52, B: 86}, c)

	c = blink.MustParseColor("#7b2d43")
	assert.Equal(t, blink.Color{R: 123, G: 45, B: 67}, c)
}
