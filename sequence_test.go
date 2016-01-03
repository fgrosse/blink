package blink_test

import (
	"time"

	"github.com/fgrosse/blink"
)

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
		Wait(1 * time.Second).
		Off()

	// blocks until s is done
	err = s.Play(led)
	if err != nil {
		panic(err)
	}
}

func ExampleSequence_Loop() {
	led, err := blink.New()
	if err != nil {
		panic(err)
	}

	defer led.FadeOutClose()

	firstLoop := blink.NewSequence().
		Fade(blink.Red, 250*time.Millisecond).
		Fade(blink.Blue, 250*time.Millisecond).
		LoopN(2) // loops 2 times

	secondLoop := firstLoop.
		Fade(blink.Green, 250*time.Millisecond).
		LoopN(4) // loops 4 times

	myBlue := blink.MustParseColor("#6666ff")
	entireLoop, c := secondLoop.
		Start(). // instruct the next loop to start at this position
		Set(myBlue, 200*time.Millisecond).
		Set(myBlue.Multiply(0.8), 200*time.Millisecond).
		Set(myBlue.Multiply(0.6), 200*time.Millisecond).
		Set(myBlue.Multiply(0.4), 200*time.Millisecond).
		Loop() // restarts the sequence until c is closed

	go func() {
		// stop the whole loop after ten seconds
		time.Sleep(10 * time.Second)
		close(c)
	}()

	err = entireLoop.Play(led)
	if err != nil {
		panic(err)
	}
}
