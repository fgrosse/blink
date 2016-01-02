package blink

import (
	"encoding/hex"
	"fmt"
	"strings"
	"strconv"

	"github.com/hashicorp/go-multierror"
)

func MustParseColor(s string) Color {
	c, err := ParseColor(s)
	if err != nil {
		panic(err)
	}

	return c
}

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
