Fire Alarm Beep Detector
====

[![fire fire fire](https://img.youtube.com/vi/1EBfxjSFAxQ/0.jpg)](https://www.youtube.com/watch?v=1EBfxjSFAxQ)

_Note: also works for water leak detector_

Featuring:

1/ An alarm beep detector
- detects fire alarm audio signature using FFT and pattern matching
- exec' a custom script if alarm detected - pagerduty & picture-upload scripts included

<br />

2/ An optional simple HTTP server that
- accepts uploads from script
- automagically serves requests with TLS (when not running locally)
- basic auth authentication
- serves last image (top level) and a the last 50 (/list)

## Motivation

I wanted a local solution running on a Raspberry Pi Zero, and _not_ constantly pushing data to a server, surely it's doable with simple FFTs and pattern matching !

## Usage

```console
$ cd detect
$ go mod tidy
$ go build
$ ./detect --help
Usage of ./detect:
  -device string
        target device. If empty, will list devices.
  -script string
        script to exec when an alarm is detected
  -threshold string
        audio target threshold (default "7")

$ ./detect -device="Logitech StreamCam: USB Audio (hw:1,0)"
[play alarm.wav]
2024-04-07 13:41:47 -- alarm detected!
```

Feel free to explore the codebase for both sources, this is more of a glorified (but fuctioning !) script :).

A pre-compiled binary for Raspberry Pi is provided on the release section. It only has a dependency on `portaudio19-dev`.
The two scripts have dependencies on `curl`, and `ffmpeg` for the picture uploading script.