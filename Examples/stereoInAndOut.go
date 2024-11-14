package main

import "github.com/cocoonlife/goalsa"

const SAMPLE_RATE = 16000

func main() {
	//hw:<CARD_NR>,<DEVICE_NR>
	bufParams := alsa.BufferParams{
		BufferFrames: 1920,
		PeriodFrames: 480,
		Periods:      4,
	}
	captureDevice, err := alsa.NewCaptureDevice("hw:2,0", 2, alsa.FormatS16LE, SAMPLE_RATE, bufParams)
	if err != nil {
		panic(err)
	}
	defer captureDevice.Close()
}
