package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"github.com/xthexder/go-jack"
)

var channels int = 2

var PortsIn []*jack.Port
var PortsOut []*jack.Port

// Global audio buffer and index
var audioBuffer []int16
var playbackIndex int = 0

func int16ToAudioSample(sample int16) jack.AudioSample {
	return jack.AudioSample(float32(sample) / 32767)
}

func AudioSampleToInt16(sample jack.AudioSample) int16 {
	clamped := math.Max(-1.0, math.Min(1.0, float64(sample)))
	return int16(clamped * 32767)
}

func process(nframes uint32) int {
	if playbackIndex >= len(audioBuffer) {
		return 0 // Stop playback when the buffer is finished
	}

	for i := 0; i < channels; i++ {
		samplesOut := PortsOut[i].GetBuffer(nframes)
		for j := 0; j < int(nframes); j++ {
			if playbackIndex < len(audioBuffer) {
				sample := audioBuffer[playbackIndex]
				samplesOut[j] = int16ToAudioSample(sample)
				playbackIndex++
			} else {
				samplesOut[j] = 0 // Fill with silence after buffer ends
			}
		}
	}

	return 0
}

func loadWavFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read the header
	header := make([]byte, 44)
	if _, err := file.Read(header); err != nil {
		return fmt.Errorf("failed to read WAV header: %v", err)
	}

	// Check the "RIFF" and "WAVE" identifiers
	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return fmt.Errorf("invalid WAV file format")
	}

	// Extract metadata
	numChannels := binary.LittleEndian.Uint16(header[22:24])
	sampleRate := binary.LittleEndian.Uint32(header[24:28])
	bitsPerSample := binary.LittleEndian.Uint16(header[34:36])

	fmt.Printf("Channels: %d, Sample Rate: %d, Bits per Sample: %d\n", numChannels, sampleRate, bitsPerSample)

	if bitsPerSample != 16 {
		return fmt.Errorf("unsupported bit depth: %d (only 16-bit PCM is supported)", bitsPerSample)
	}
	if int(numChannels) != channels {
		return fmt.Errorf("expected %d channels, but found %d", channels, numChannels)
	}

	// Read the audio data
	audioData := make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		audioData = append(audioData, buffer[:n]...)
	}

	// Convert byte buffer to int16 samples
	audioBuffer = make([]int16, len(audioData)/2)
	for i := 0; i < len(audioBuffer); i++ {
		audioBuffer[i] = int16(binary.LittleEndian.Uint16(audioData[i*2 : i*2+2]))
	}

	fmt.Printf("Loaded %d samples\n", len(audioBuffer))
	return nil
}

func main() {
	myClient, status := jack.ClientOpen("Go Passthrough", jack.NoStartServer)
	if status != 0 {
		fmt.Println("Status1:", jack.StrError(status))
		return
	}
	defer myClient.Close()

	portOut1 := myClient.PortRegister("Lane_1_Output", jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
	portOut2 := myClient.PortRegister("Lane_2_Output", jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
	PortsOut = append(PortsOut, portOut1, portOut2)

	// Load the WAV file
	err := loadWavFile("../CantinaBand3.wav")
	if err != nil {
		fmt.Println("Error loading WAV file:", err)
		return
	}

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
