package blink

import "time"

const reportID = 0x01

type command interface {
	bytes() []byte
}

type setRGBCommand struct {
	Color // 24-bit RGB color
}

func (c *setRGBCommand) bytes() []byte {
	return []byte{reportID,
		'n',
		c.r, c.g, c.b,
		0, 0, 0,
	}
}

type fadeRGBCommand struct {
	Color                  // 24-bit RGB color
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

type readRGBCommand struct {
	n byte // which LED to address: 0=all, 1=led#1, 2=led#2, etc. (mk2 only)
}

func (c *readRGBCommand) bytes() []byte {
	return []byte{reportID,
		'r',
		0, 0, 0,
		0, 0,
		c.n,
	}
}
