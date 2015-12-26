package blink

import "time"

const reportID = 0x01

type command interface {
	bytes() []byte
}

type rgb struct{ r, g, b byte }

type setRGBCommand struct {
	rgb // 24-bit RGB color
}

func (c *setRGBCommand) bytes() []byte {
	return []byte{reportID,
		'n',
		c.r, c.g, c.b,
		0, 0, 0,
	}
}

type fadeRGBCommand struct {
	rgb                    // 24-bit RGB color
	duration time.Duration // how long the fade should last
	n        byte          // which LED to address: 0=all, 1=led#1, 2=led#2, etc. (mk2 only)
}

func (c *fadeRGBCommand) bytes() []byte {
	t := c.duration.Nanoseconds() / 1E7
	return []byte{reportID,
		'c',
		c.r, c.g, c.b,
		byte(t >> 8), byte(t & 0xff),
		c.n,
	}
}
