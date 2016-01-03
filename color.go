package blink

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
)

var (
	Red    = Color{R: 255, G: 000, B: 000}
	Green  = Color{R: 000, G: 255, B: 000}
	Blue   = Color{R: 000, G: 000, B: 255}
	Yellow = Color{R: 255, G: 255, B: 000}
	White  = Color{R: 255, G: 255, B: 255}
)

// Color contains the 24-bit RGB color information.
type Color struct{ R, G, B byte }

// Multiply returns a copy of c where every individual color channel is
// multiplied with the given factor.
// If f is lower than 0 it is ignored entirely and c is returned unchanged.
func (c Color) Multiply(f float64) Color {
	if f < 0 {
		return c
	}

	return Color{
		R: floatToByte(float64(c.R) * f),
		G: floatToByte(float64(c.G) * f),
		B: floatToByte(float64(c.B) * f),
	}
}

func floatToByte(f float64) byte {
	if f > 255 {
		f = 255
	}

	return byte(f + 0.5)
}

// MustParseColor behaves exactly as ParseColor but panics if an error occurs.
func MustParseColor(s string) Color {
	c, err := ParseColor(s)
	if err != nil {
		panic(err)
	}

	return c
}

// ParseColor parses a Color from string.
// It accepts the RGB value either as comma separated vector
// or in hexadecimal form with a leading hash tag.
// Examples:
//     255,255,0
//     255, 255, 0
//     #ffff00
//     #FFFF00
func ParseColor(s string) (c Color, err error) {
	switch {
	case len(s) == 0:
		err = fmt.Errorf("can not parse color from empty string")
	case s[0] == '#':
		c, err = parseColorFromHex(s[1:])
	default:
		c, err = parseColorFromCSV(s)
	}

	return
}

func parseColorFromHex(s string) (Color, error) {
	var c Color
	b, err := hex.DecodeString(s)
	if err != nil {
		return c, fmt.Errorf("can not parse hex color from %q: %s", s, err)
	}

	if len(b) != 3 {
		return c, fmt.Errorf("invalid number of bytes (have %d want 3)", len(b))
	}

	c.R = b[0]
	c.G = b[1]
	c.B = b[2]

	return c, nil
}

func parseColorFromCSV(s string) (Color, error) {
	var c Color
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return c, fmt.Errorf("can not parse color from CSV: expected exactly three comma separated values")
	}

	parse := func(s string) (byte, error) {
		i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 8)
		return byte(i), err
	}

	var err1, err2, err3 error
	c.R, err1 = parse(parts[0])
	c.G, err2 = parse(parts[1])
	c.B, err3 = parse(parts[2])

	if err1 != nil || err2 != nil || err3 != nil {
		return c, multierror.Append(err1, err2, err3)
	}

	return c, nil
}
