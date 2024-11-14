package main

import (
	"fmt"

	"github.com/xthexder/go-jack"
)

var channels int = 2

var PortsIn []*jack.Port
var PortsOut []*jack.Port

func process(nframes uint32) int {
	count := 0
	for i, in := range PortsIn {
		samplesIn := in.GetBuffer(nframes)
		samplesOut := PortsOut[i].GetBuffer(nframes)
		for i2, sample := range samplesIn {
			samplesOut[i2] = sample
			count += 1
		}
	}
	fmt.Println("Samples Processed: %v", count)
	return 0
}

func main() {
	myClient, status := jack.ClientOpen("Go Passthrough", jack.NoStartServer)
	if status != 0 {
		fmt.Println("Status1:", jack.StrError(status))
		return
	}
	defer myClient.Close()

	echoCancelSource, status := jack.ClientOpen("Echo-Cancel Source", jack.NoStartServer)
	if status != 0 {
		fmt.Println("Status2:", jack.StrError(status))
		return
	}
	defer echoCancelSource.Close()

	echoCancelSink, status := jack.ClientOpen("Echo-Cancel Sink", jack.NoStartServer)
	if status != 0 {
		fmt.Println("Status3:", jack.StrError(status))
		return
	}
	defer echoCancelSink.Close()

	portIn1 := myClient.PortRegister("Lane1Input", jack.DEFAULT_AUDIO_TYPE, jack.PortIsInput, 0)
	portIn2 := myClient.PortRegister("Lane2Input", jack.DEFAULT_AUDIO_TYPE, jack.PortIsInput, 0)
	PortsIn = append(PortsIn, portIn1, portIn2)
	myClient.ConnectPorts(echoCancelSource.GetPortByName("capture_FL"), portIn1)
	myClient.ConnectPorts(echoCancelSource.GetPortByName("capture_FR"), portIn2)

	portOut1 := myClient.PortRegister("Lane_1_Output", jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
	portOut2 := myClient.PortRegister("Lane_2_Output", jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
	PortsOut = append(PortsOut, portOut1, portOut2)
	myClient.ConnectPorts(echoCancelSink.GetPortByName("playback_FL"), portIn1)
	myClient.ConnectPorts(echoCancelSink.GetPortByName("playback_FR"), portIn2)

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
