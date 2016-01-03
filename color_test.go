package blink_test

import (
	"testing"

	"fmt"
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

func TestMultiply(t *testing.T) {
	data := []struct {
		have blink.Color
		m    float64
		want blink.Color
	}{
		{blink.Color{R: 124, G: 64, B: 246}, 0.0, blink.Color{R: 0, G: 0, B: 0}},
		{blink.Color{R: 124, G: 64, B: 246}, 0.5, blink.Color{R: 62, G: 32, B: 123}},
		{blink.Color{R: 1, G: 73, B: 255}, 0.7, blink.Color{R: 1, G: 51, B: 179}},
		{blink.Color{R: 124, G: 64, B: 246}, 1.0, blink.Color{R: 124, G: 64, B: 246}},
		{blink.Color{R: 1, G: 73, B: 255}, 1.5, blink.Color{R: 2, G: 110, B: 255}},
		{blink.Color{R: 1, G: 73, B: 255}, 2.0, blink.Color{R: 2, G: 146, B: 255}},
		{blink.Color{R: 1, G: 73, B: 255}, -1, blink.Color{R: 1, G: 73, B: 255}}, // should be ignored
	}

	for i, d := range data {
		assert.Equal(t, d.want, d.have.Multiply(d.m), fmt.Sprintf("%d: %+v * %.1f", i, d.have, d.m))
	}
}

func TestAdd(t *testing.T) {
	data := []struct{ have, add, want blink.Color }{
		{
			blink.Color{R: 0, G: 0, B: 0},
			blink.Color{R: 0, G: 0, B: 0},
			blink.Color{R: 0, G: 0, B: 0},
		},
		{
			blink.Color{R: 10, G: 20, B: 30},
			blink.Color{R: 0, G: 0, B: 0},
			blink.Color{R: 10, G: 20, B: 30},
		},
		{
			blink.Color{R: 10, G: 40, B: 100},
			blink.Color{R: 20, G: 50, B: 200},
			blink.Color{R: 30, G: 90, B: 255},
		},
	}

	for i, d := range data {
		assert.Equal(t, d.want, d.have.Add(d.add), fmt.Sprintf("%d: %+v + %+v", i, d.have, d.add))
	}
}
