package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Binozo/GoAlsa/pkg/alsa"
)

const SAMPLE_RATE = 16000
const BUFFER_SIZE = 480

func main() {
	//hw:<CARD_NR>,<DEVICE_NR>
	// bufParams := alsa.BufferParams{
	// 	BufferFrames: 1920,
	// 	PeriodFrames: 480,
	// 	Periods:      4,
	// }
	audioConfig := alsa.Config{
		Channels:   2,
		Format:     alsa.FormatS16LE,
		SampleRate: SAMPLE_RATE,
	}
	captureDevice, err := alsa.NewCaptureDevice("hw:2,0", audioConfig)
	if err != nil {
		panic(err)
	}
	defer captureDevice.Close()

	playbackDevice, err := alsa.NewPlaybackDevice("hw:2,0", audioConfig)
	if err != nil {
		panic(err)
	}
	defer playbackDevice.Close()

	//captureDevice.StartReadThread()

	readBuffer := make([]float32, 480*2)
	writeBuffer := make([]float32, 480*2)

	go func() (err error) {
		for {
			numSamples, err := captureDevice.Read(readBuffer)
			if err != nil {
				return fmt.Errorf("error reading capture device, %v", err)
			}
			fmt.Println("Num samples in last read: ", numSamples)
			copy(writeBuffer, readBuffer)
			numSamples, err = playbackDevice.Write(writeBuffer)
			if err != nil {
				return fmt.Errorf("error writing to playback device, %v", err)
			}
			fmt.Println("Num samples in last write: ", numSamples)
			time.Sleep(time.Millisecond * 20)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	fmt.Println("Exiting...")
}
