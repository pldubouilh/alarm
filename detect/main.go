package main

import (
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/mjibson/go-dsp/fft"
)

const (
	sampleRate  = 48000
	frameSize   = 512
	numChannels = 1
)

var device = flag.String("device", "", "Target device. If empty, will list devices.")
var script = flag.String("script", "", "Script to exec when an alarm is detected")
var threshold_str = flag.String("threshold", "7", "Audio target threshold")
var threshold float64

var targetFreq = flag.Int("frequency", 3500, "Target frequency in Hz")
var beepDuration = flag.Duration("duration", 400*time.Millisecond, "Duration of a beep")
var howManyBeepsToAlert = flag.Int("beeps", 3, "How many beeps to aloers")

func main() {
	flag.Parse()
	i, err := strconv.Atoi(*threshold_str)
	if err != nil {
		fmt.Println("Failed to parse threshold:", err)
		os.Exit(1)
	}
	threshold = float64(i) * 1_000_000_000

	for {
		if !run() {
			println("Failed to find device, will retry in 5 seconds")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
}

func run() bool {
	// Initialize PortAudio
	if err := portaudio.Initialize(); err != nil {
		fmt.Println("Failed to initialize PortAudio:", err)
		return false
	}
	defer portaudio.Terminate()

	var info *portaudio.DeviceInfo
	ds, err := portaudio.Devices()
	if err != nil {
		fmt.Println("Failed to get devices:", err)
		return false
	}
	for i, d := range ds {
		if *device == "" {
			fmt.Printf("Found device %s\n", d.Name)
		}
		if d.Name == *device {
			info = ds[i]
			break
		}
	}

	if *device == "" {
		os.Exit(0)
	}
	if info == nil {
		fmt.Println("Failed to find device:", *device)
		return false
	}

	// Prepare input parameters for the stream
	inputParams := portaudio.StreamParameters{
		Input:           portaudio.StreamDeviceParameters{Device: info, Channels: numChannels, Latency: info.DefaultHighInputLatency},
		Output:          portaudio.StreamDeviceParameters{Device: nil, Channels: 0, Latency: 0},
		SampleRate:      sampleRate,
		FramesPerBuffer: frameSize,
	}

	// Open stream
	stream, err := portaudio.OpenStream(inputParams, processAudio)
	if err != nil {
		fmt.Println("Failed to open stream:", err)
		return false
	}
	defer stream.Close()

	// Start stream
	if err := stream.Start(); err != nil {
		fmt.Println("Failed to start stream:", err)
		return false
	}
	defer stream.Stop()

	fmt.Println("Listening for audio...")

	// Check beeps and bops
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
	fmt.Println("\nExiting...")
	return true
}

func processAudio(in []int32) {
	// Windowing function to the audio samples before performing the FFT.
	window := make([]float64, len(in))
	for i, x := range in {
		window[i] = float64(x) * (0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(len(in)-1)))
	}

	// Direct FFT
	// window := make([]float64, len(in))
	// for i, x := range in {
	// 	window[i] = float64(x)
	// }

	fftData := fft.FFTReal(window)

	// Find the magnitude of the target frequency bin
	targetIndex := int(float64(len(fftData)) * float64(*targetFreq) / float64(sampleRate))
	magnitude := cmplx.Abs(fftData[targetIndex])

	// Check if the magnitude is above the threshold
	if magnitude > threshold {
		go checkBeeps(true)
		// fmt.Printf("Significant data detected %dHz, magnitude: %.0f\n", targetFreq, magnitude)
	} else {
		go checkBeeps(false)
		// fmt.Printf("\rNo significant magnitude at %d Hz\n", targetFreq)
	}
}

var lastChange = time.Now()
var beeps = 0
var current = false

func checkBeeps(state bool) {
	if state == current {
		return
	}
	current = state
	delta := time.Since(lastChange)
	lastChange = time.Now()

	// too old
	if delta > 2*(*beepDuration) {
		beeps = 0
		return
	}

	if !state && delta > *beepDuration {
		beeps++
	}

	if beeps >= *howManyBeepsToAlert {
		go execScript()
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "-- alarm detected!")
		beeps = 0
	}
}

func execScript() {
	if *script == "" {
		return
	}
	cmd := exec.Command("/bin/sh", "-c", *script)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to run execScript script:", err)
	}
}
