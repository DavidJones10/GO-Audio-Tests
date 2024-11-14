package main

import (
	"fmt"
	"math"

	"github.com/xthexder/go-jack"
)

var channels int = 2

var PortsIn []*jack.Port
var PortsOut []*jack.Port

func int16ToAudioSample(sample int16) jack.AudioSample {
	return jack.AudioSample(float32(sample) / 32767)
}

func AudioSampleToInt16(sample jack.AudioSample) int16 {
	clamped := math.Max(-1.0, math.Min(1.0, float64(sample)))
	return int16(clamped * 32767)
}

var maxFloat float32 = 0

func process(nframes uint32) int {
	for i, in := range PortsIn {
		samplesIn := in.GetBuffer(nframes)
		samplesOut := PortsOut[i].GetBuffer(nframes)
		for i2, sample := range samplesIn {
			intSample := AudioSampleToInt16(sample)
			samplesOut[i2] = int16ToAudioSample(intSample)
			fmt.Println("Sample: ", intSample)
			if samplesOut[i2] > jack.AudioSample(maxFloat) {
				fmt.Println("Max Float: ", maxFloat)
			}
		}
	}
	return 0
}

func main() {
	myClient, status := jack.ClientOpen("Go Passthrough", jack.NoStartServer)
	if status != 0 {
		fmt.Println("Status1:", jack.StrError(status))
		return
	}
	defer myClient.Close()

	portIn1 := myClient.PortRegister("Lane1Input", jack.DEFAULT_AUDIO_TYPE, jack.PortIsInput, 0)
	portIn2 := myClient.PortRegister("Lane2Input", jack.DEFAULT_AUDIO_TYPE, jack.PortIsInput, 0)
	PortsIn = append(PortsIn, portIn1, portIn2)

	portOut1 := myClient.PortRegister("Lane_1_Output", jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
	portOut2 := myClient.PortRegister("Lane_2_Output", jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
	PortsOut = append(PortsOut, portOut1, portOut2)

	if code := myClient.SetProcessCallback(process); code != 0 {
		fmt.Println("Failed to set process callback:", jack.StrError(code))
		return
	}
	shutdown := make(chan struct{})
	myClient.OnShutdown(func() {
		fmt.Println("Shutting down")
		close(shutdown)
	})

	if code := myClient.Activate(); code != 0 {
		fmt.Println("Failed to activate client:", jack.StrError(code))
		return
	}

	fmt.Println(myClient.GetName())
	<-shutdown
}
