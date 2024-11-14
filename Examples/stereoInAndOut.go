package main

import (
	"log"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100
const bufferSize = 2048
const numChannels = 2

func main() {
	// Initialize PortAudio
	if err := portaudio.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer portaudio.Terminate()

	// Create an input buffer to store the incoming audio data
	// inputBuffer := make([]int16, bufferSize*numChannels)
	// outputBuffer := make([]int16, bufferSize*numChannels)

	// Open the default input and output streams
	stream, err := portaudio.OpenDefaultStream(numChannels, numChannels, sampleRate, bufferSize, processAudio)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	// Start the stream
	if err := stream.Start(); err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()
}

func processAudio(in, out []int16) {
	for sample := 0; sample < len(out); sample++ {
		out[sample] = in[sample]
	}
}
