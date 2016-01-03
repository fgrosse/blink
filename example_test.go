package blink_test

import (
	"fmt"
	"time"

	"github.com/fgrosse/blink"
)

func Example() {
	// connect to a local blink(1) USB device
	led, err := blink.New()
	if err != nil {
		panic(err)
	}

	// disable all lights and close the device when you are done
	defer led.FadeOutClose()

	// fade to a full green in 500ms
	d := 500 * time.Millisecond
	led.FadeRGB(0, 255, 0, d)
	time.Sleep(d)

	// store colors for later use
	corpBlue := blink.MustParseColor("#3333ff")
	led.Fade(corpBlue, d)

	// read the current color
	color, err := led.Read()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", color)
}
