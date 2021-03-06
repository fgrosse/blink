# blink

[![Build Status](https://travis-ci.org/fgrosse/blink.svg?branch=master)](https://travis-ci.org/fgrosse/blink)
[![GoDoc](https://godoc.org/github.com/fgrosse/blink?status.svg)](https://godoc.org/github.com/fgrosse/blink)
[![License](https://img.shields.io/badge/license-MIT-4183c4.svg)](https://github.com/fgrosse/blink/blob/master/LICENSE)

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

## Installation

Currently blink does only compile on **linux** and requires **[libusb-1.0.12][5] or higher**.
blink is build on travis using libusb 1.0.20. Refer to the [`.travis.yml`](.travis.yml) to it can be built on ubuntu.
On Fedora 22 you can simply use `dnf install libusb-devel`.

Use `go get` to install blink:
```
go get github.com/fgrosse/blink
```

You need to have go version 1.4 or higher.

## Usage

Simple usage
```go
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
```

Create **sequences** to store and playback multiple instructions
```go
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
```

Sequences can be run in a loop. You can also loop multiple sections.
```go
firstLoop := blink.NewSequence().
    Fade(blink.Red, 250*time.Millisecond).
    Fade(blink.Blue, 250*time.Millisecond).
    LoopN(2) // loops 2 times

secondLoop := firstLoop.
    Fade(blink.Green, 250*time.Millisecond).
    LoopN(4) // loops 4 times starting at the first fade to red

myBlue := blink.MustParseColor("#6666ff")
entireLoop, c := secondLoop.
    Start(). // instruct the next loop to start at this position
    Set(myBlue, 200 * time.Millisecond).
    Set(myBlue.Multiply(0.8), 200 * time.Millisecond).
    Set(myBlue.Multiply(0.6), 200 * time.Millisecond).
    Set(myBlue.Multiply(0.4), 200 * time.Millisecond).
    Loop() // restarts the sequence until c is closed

go func() {
    // stop the whole loop after ten seconds
    time.Sleep(10 * time.Second)
    close(c)
}()

err = entireLoop.Play(led)
```
### Linux Permissions

You need to have root access when running this program or you will get the following error:

```
libusb: bad access [code -3]
```

On linux this problem can easily be fixed by adding the following [udev rule][6]:

```bash
[root@localhost]# cat /etc/udev/rules.d/10.local.rules
SUBSYSTEMS=="usb", ATTRS{idVendor}=="27b8", ATTRS{idProduct}=="01ed", SYMLINK+="blink1", GROUP="blink1"
```

Everybody in the `blink1` group should now be able to access the device directly.
Additionally this rule creates a symlink at `/dev/blink1` each time you connect the device.
You probably need to reconnect your device so the change will be visible.

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
[5]: https://github.com/libusb/libusb
[6]: http://www.reactivated.net/writing_udev_rules.html
