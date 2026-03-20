package waveform

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/wav"
)

func AnalyzeFile(path string) ([][2]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	streamer := beep.StreamSeekCloser(nil)

	if filepath.Ext(path) == ".mp3" {
		streamer, _, err = mp3.Decode(f)
	} else if filepath.Ext(path) == ".wav" {
		streamer, _, err = wav.Decode(f)
	} else if filepath.Ext(path) == ".flac" {
		streamer, _, err = flac.Decode(f)
	} else {
		return nil, fmt.Errorf("unsupported file format: %s", filepath.Ext(path))
	}
	if err != nil {
		fmt.Printf("Error decoding file %s: %v\n", path, err)
		return nil, err
	}
	defer streamer.Close()

	totalSamples := streamer.Len()
	if totalSamples == 0 {
		return nil, fmt.Errorf("no audio data found in file: %s", path)
	}

	samplesPerPoint := totalSamples / 10000
	if samplesPerPoint == 0 {
		samplesPerPoint = 1
	}

	waveform := make([][2]float64, 0, 10000)
	buffer := make([][2]float64, samplesPerPoint)

	for i := 0; i < 10000; i++ {
		n, ok := streamer.Stream(buffer)
		if !ok {
			break
		}
		fmt.Printf("Processing point %d: read %d samples\n", i, n)
		maxAmplitude := [2]float64{0.0, 0.0}

		for j := 0; j < n; j++ {
			left := math.Abs(buffer[j][0])
			right := math.Abs(buffer[j][1])
			if left > maxAmplitude[0] {
				maxAmplitude[0] = left
			}
			if right > maxAmplitude[1] {
				maxAmplitude[1] = right
			}
		}
		waveform = append(waveform, maxAmplitude)
	}
	fmt.Printf("Analyzed %d samples, generated %d waveform points\n", totalSamples, len(waveform))
	return waveform, nil
}
