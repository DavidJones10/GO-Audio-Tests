package main

import (
	"log"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100
const bufferSize = 512
const numChannels = 2

func main() {
	// Initialize PortAudio
	if err := portaudio.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer portaudio.Terminate()

	// Create an input buffer to store the incoming audio data
	inputBuffer := make([]int16, bufferSize*numChannels)
	outputBuffer := make([]int16, bufferSize*numChannels)

	// Open the default input and output streams
	stream, err := portaudio.OpenDefaultStream(numChannels, numChannels, sampleRate, bufferSize, inputBuffer)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	// Start the stream
	if err := stream.Start(); err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()

	// Main loop to read input and immediately write it to the output
	for {
		// Read input audio
		if err := stream.Read(); err != nil {
			log.Fatal(err)
		}

		// Copy the input to the output buffer (loopback)
		copy(outputBuffer, inputBuffer)

		// Write output audio
		if err := stream.Write(); err != nil {
			log.Fatal(err)
		}
	}
}
