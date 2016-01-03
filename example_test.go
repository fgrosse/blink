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

func ExampleSequence() {
	led, err := blink.New()
	if err != nil {
		panic(err)
	}

	defer led.Close()

	d := 500 * time.Millisecond
	s := blink.NewSequence().
		Fade(blink.Red, d).
		Fade(blink.Green, d).
		Fade(blink.Blue, d).
		Off()

	// blocks until s is done
	err = s.Play(led)
	if err != nil {
		panic(err)
	}
}

func ExampleSequenceLoop() {
	led, err := blink.New()
	if err != nil {
		panic(err)
	}

	defer led.FadeOutClose()

	d := 500 * time.Millisecond
	police, c := blink.NewSequence().
		Fade(blink.Red, d).
		Fade(blink.Blue, d).
		Loop()

	go func() {
		time.Sleep(3 * 2 * d)
		close(c)
	}()

	err = police.Play(led)
	if err != nil {
		panic(err)
	}
}
