package blink

import (
	"fmt"
	"time"
)

type Sequence struct {
	frames []frame
}

type frame interface {
	run(int, *LED) (int, error)
}

func NewSequence() *Sequence {
	return &Sequence{}
}

type cmdFrame struct {
	command
	time.Duration
}

func (f *cmdFrame) run(step int, led *LED) (int, error) {
	_, err := led.write(f.command)
	if err != nil {
		return 0, err
	}

	time.Sleep(f.Duration)
	return step+1, nil
}

type loopFrame struct {
	stop bool
}

func (f *loopFrame) run(step int, led *LED) (int, error) {
	if f.stop {
		return step+1, nil
	}

	return 0, nil
}

func (s *Sequence) Off() *Sequence {
	s.frames = append(s.frames, &cmdFrame{
		command: &setRGBCommand{Color{}},
	})

	return s
}

func (s *Sequence) Fade(c Color, d time.Duration) *Sequence {
	s.frames = append(s.frames, &cmdFrame{
		command:  &fadeRGBCommand{Color: c, duration: d},
		Duration: d,
	})

	return s
}

func (s *Sequence) Loop() (*Sequence, chan<- bool) {
	f := new(loopFrame)
	s.frames = append(s.frames, f)

	c := make(chan bool)
	go func() {
		for range c {}
		f.stop = true
	}()

	return s, c
}

func (s *Sequence) Play(led *LED) error {
	if led == nil {
		return fmt.Errorf("led is nil")
	}

	var i int
	var err error
	for {
		if i >= len(s.frames) {
			break
		}

		if i, err = s.frames[i].run(i, led); err != nil {
			return err
		}
	}

	return nil
}
