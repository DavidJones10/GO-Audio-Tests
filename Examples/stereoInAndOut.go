package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/cocoonlife/goalsa"
)

const SAMPLE_RATE = 16000
const BUFFER_SIZE = 480

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

	captureDevice.StartReadThread()
	readBuffer := make([]int16, 480*2)

	go func() (err error) {
		for {
			numSamples, err := captureDevice.Read(readBuffer)
			fmt.Println("Num samples in last read: ", numSamples)
			if err != nil {
				return fmt.Errorf("error reading capture device")
			}

		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	fmt.Println("Exiting...")
}
