package main

import (
	"github.com/gordonklaus/portaudio"
)

const BUFFER_SIZE = 480
const SAMPLE_RATE = 16000

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	inBuf := make([][]int16, BUFFER_SIZE)
	inStream, err := portaudio.OpenDefaultStream(1, 0, SAMPLE_RATE, BUFFER_SIZE, inBuf)
	if err != nil {
		panic(err)
	}

	outBuf := make([][]int16, BUFFER_SIZE)
	outStream, err := portaudio.OpenDefaultStream(0, 1, SAMPLE_RATE, BUFFER_SIZE, outBuf)
	if err != nil {
		panic(err)
	}

	inStream.Start()
	outStream.Start()

	go processAudio(inStream, outStream, inBuf, outBuf)

}

func processAudio(inputStream *portaudio.Stream, outputStream *portaudio.Stream, inputBuffer [][]int16, outputBuffer [][]int16) {
	for {
		if err := inputStream.Read(); err != nil {
			panic(err)
		}

		for channel := 0; channel < 1; channel++ {
			for i := 0; i < BUFFER_SIZE; i++ {
				outputBuffer[channel][i] = inputBuffer[channel][i] // Loop input to output
			}
		}
		if err := outputStream.Write(); err != nil {
			panic(err)
		}
	}
}
