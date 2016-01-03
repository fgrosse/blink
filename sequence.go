package blink

import (
	"fmt"
	"time"
)

// A Sequence is used to store and send a set of commands to the blink(1) device
// in particular order. It enables you to define a sequence of actions like
// fading to another color, waiting a certain amount of time and looping.
//
// A Sequence is not save for concurrent use.
type Sequence struct {
	frames []frame
	i      int // the current frame
}

type frame interface {
	run(*LED) error
}

// NewSequence creates a new sequence and can be used to chain multiple sequence instructions
// Example:
//     s := blink.NewSequence().
//         Fade(blink.Red, d).
//         Fade(blink.Green, d).
//         Fade(blink.Blue, d).
//         Wait(1 * time.Second).
//         Off()
func NewSequence() *Sequence {
	return &Sequence{}
}

// Off adds a new frame to the sequence which deactivates the led immediately.
func (s *Sequence) Off() *Sequence {
	return s.Set(Color{}, 0)
}

// Set adds a new frame to the sequence which immediately sets the led to another
// color and waits a given duration.
func (s *Sequence) Set(c Color, d time.Duration) *Sequence {
	s.frames = append(s.frames, &cmdFrame{
		command:  &setRGBCommand{Color: c},
		Duration: d,
	})
	return s
}

// Fade adds a new frame to the sequence which lets the led fade to another color.
func (s *Sequence) Fade(c Color, d time.Duration) *Sequence {
	s.frames = append(s.frames, &cmdFrame{
		command:  &fadeRGBCommand{Color: c, duration: d},
		Duration: d,
	})

	return s
}

// Wait adds a new frame to the sequence which doesn't do anything for a given duration.
func (s *Sequence) Wait(d time.Duration) *Sequence {
	s.frames = append(s.frames, &waitFrame{d})
	return s
}

// Loop is used to instruct the sequence to loop all previously added frames infinitely.
// The second return value is a channel that can be closed to stop the looping sequence
// after it has been started. Sending values to the channel does nothing.
// The loop frame works by resetting the sequence each time it is reached.
// Once the channel has been closed the loop will stop to reset the sequence.
func (s *Sequence) Loop() (*Sequence, chan<- struct{}) {
	f := &loopFrame{seq: s, n: -1}
	s.frames = append(s.frames, f)

	c := make(chan struct{})
	go func() {
		for range c {
			// just wait for c to be closed
		}
		f.n = 0
	}()

	return s, c
}

func (s *Sequence) LoopN(n int) *Sequence {
	s.frames = append(s.frames, &loopFrame{seq: s, n: n})
	return s
}

func (s *Sequence) Start() *Sequence {
	f := &startFrame{seq: s, n: len(s.frames)}
	s.frames = append(s.frames, f)
	return s
}

// Play starts to playback this sequence on the given LED.
// It blocks until all frames have been processed.
// If this sequence loops Play will never return by itself.
func (s *Sequence) Play(led *LED) error {
	if led == nil {
		return fmt.Errorf("led is nil")
	}

	s.i = 0
	var err error
	for {
		if s.i >= len(s.frames) {
			break
		}

		f := s.frames[s.i]
		if err = f.run(led); err != nil {
			return err
		}

		s.i++
	}

	return nil
}

type cmdFrame struct {
	command
	time.Duration
}

func (f *cmdFrame) run(led *LED) error {
	_, err := led.write(f.command)
	if err != nil {
		return err
	}

	time.Sleep(f.Duration)
	return nil
}

type waitFrame struct{ time.Duration }

func (f *waitFrame) run(led *LED) error {
	time.Sleep(f.Duration)
	return nil
}

type loopFrame struct {
	seq *Sequence
	n   int
}

func (f *loopFrame) run(led *LED) error {
	if f.n > 0 {
		f.n--
	}

	if f.n != 0 {
		f.seq.i = -1 // works because this is always incremented after each frame
	}

	return nil
}

type startFrame struct {
	seq *Sequence
	n   int
}

func (f *startFrame) run(led *LED) error {
	f.seq.frames = f.seq.frames[f.n+1:]
	f.seq.i = 0
	return nil
}
