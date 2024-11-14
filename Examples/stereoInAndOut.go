package main

import (
	"fmt"
	"time"

	"github.com/gordonklaus/portaudio"
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()
	e := newEcho(time.Second / 30)
	defer e.Close()
	chk(e.Start())
	time.Sleep(4 * time.Second)
	chk(e.Stop())
}

type echo struct {
	*portaudio.Stream
	buffer []float32
	i      int
}

func newEcho(delay time.Duration) *echo {
	h, err := portaudio.DefaultHostApi()
	chk(err)
	fmt.Println(h.DefaultInputDevice.MaxInputChannels)
	fmt.Println(h.DefaultOutputDevice.MaxOutputChannels)
	fmt.Println(h.Name)
	fmt.Println(h.DefaultInputDevice.Name)
	fmt.Println(h.DefaultOutputDevice.Name)

	h.DefaultInputDevice.Name = "hw:2,0"
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, h.DefaultOutputDevice)
	p.Input.Channels = 1
	p.Output.Channels = 1
	fmt.Println("made it 1")
	e := &echo{buffer: make([]float32, int(p.SampleRate*delay.Seconds()))}
	e.Stream, err = portaudio.OpenStream(p, e.processAudio)
	fmt.Println("made it 2")
	chk(err)
	return e
}

func (e *echo) processAudio(in, out []float32) {
	for i := range out {
		out[i] = .7 * e.buffer[e.i]
		e.buffer[e.i] = in[i]
		e.i = (e.i + 1) % len(e.buffer)
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
