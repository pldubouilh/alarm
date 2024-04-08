Fire Alarm Detector
====

[![fire fire fire](https://img.youtube.com/vi/1EBfxjSFAxQ/0.jpg)](https://www.youtube.com/watch?v=1EBfxjSFAxQ)

_Note: also works for water leak detector_

Featuring:
- fire alarm audio signature detection using FFT and pattern matching
- exec' a custom script if the alarm is detected - example pagerduty & picture-snap-n-upload scripts provided
- an optional server to handle the uploaded pictures (with automatic https and basic auth)

## Motivation

I wanted a local solution running on a Raspberry Pi Zero, and _not_ constantly pushing data to a server, surely it's doable with simple FFTs and pattern matching !

Optionally, being able to alert with pagerduty, or take pictures when the alarm is detected (server code included).

## Usage

```console
$ cd detect
$ go mod tidy
$ go build
$ ./detect --help
Usage of ./detect:
  -beeps int
        How many beeps to alert (default 3)
  -device string
        Target device. If empty, will list devices.
  -duration duration
        Duration of a beep (default 400ms)
  -frequency int
        Target frequency in Hz (default 3500)
  -script string
        Script to exec when an alarm is detected
  -threshold string
        Audio target threshold (default "7")

$ ./detect -device="Logitech StreamCam: USB Audio (hw:1,0)"
[play alarm.wav]
2024-04-07 13:41:47 -- alarm detected!
```

Feel free to explore the codebase, this is more of a glorified (but fuctioning !) script :).

A pre-compiled binary for Raspberry Pi is provided in the release section. It only has a dependency on `portaudio19-dev`.
The two scripts have dependencies on `curl`, and `ffmpeg` for the picture uploading script.