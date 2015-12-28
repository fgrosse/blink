# blink

[![Build Status](https://travis-ci.org/fgrosse/blink.svg?branch=master)](https://travis-ci.org/fgrosse/blink)
[![GoDoc](https://godoc.org/github.com/fgrosse/blink?status.svg)](https://godoc.org/github.com/fgrosse/blink)

blink is a go implementation for controlling [ThingM blink(1) USB dual RGB LEDs][1].

## Features

- [x] Fade to RGB color
- [x] Set RGB color now  
- [x] Read current RGB color *(mk2 devices only)*
- [ ] Serverdown tickle/off
- [ ] Play/Pause
- [ ] PlayLoop *(mk2 devices only)*
- [ ] Playstate readback *(mk2 devices only)*
- [ ] Set color pattern line
- [ ] read color pattern line
- [ ] Save color patterns *(mk2 devices only)*
- [ ] Read EEPROM location *(mk1 devices only)*
- [ ] Write EEPROM location *(mk1 devices only)*
- [ ] Get version
- [ ] Test command

Eventually all of the [available HID commands will be implemented][2]

## Installation

Currently blink does only compile on linux and **requires libusb-1.0**.

Use `go get` to install blink:
```
go get github.com/fgrosse/blink
```

## Usage

```go
// connect to a local blink(1) USB device
led, err := blink.New()
if err != nil {
    panic(err)
}

// make sure its closed when you are done
defer led.Close()

// fade to a full green in 500ms
d := 500 * time.Millisecond
led.FadeRGB(0, 255, 0, d)
time.Sleep(d)

// read the current color
color, err := led.ReadRGB()
if err != nil {
    panic(err)
}
fmt.Printf("%#v\n", color)

// immediately set the color (0, 0, 0 effectively disables the led)
err = led.SetRGB(0, 0, 0)
```

## Other resources

* [the official ThingM/blink1 github repository with APIs for other languages][3]
* [GoBlink by ThingM (apparently tested on Mac OSX Mountain Lion)][4]

## Contributing

Any contributions are always welcome (use pull requests).
Please keep in mind that I might not always be able to respond immediately but I usually try to react within the week ☺.

[1]: http://blink1.thingm.com/
[2]: https://github.com/ThingM/blink1/blob/master/docs/blink1-hid-commands.md
[3]: https://github.com/ThingM/blink1
[4]: https://github.com/ThingM/blink1/tree/master/go/GoBlink
