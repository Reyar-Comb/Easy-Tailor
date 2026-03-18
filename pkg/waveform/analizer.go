package waveform

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/gopxl/beep/v2/mp3"
)

func AnalyzeFile(path string) ([][2]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if filepath.Ext(path) != ".mp3" {
		return nil, fmt.Errorf("unsupported file format: %s", filepath.Ext(path))
	}
	streamer, _, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	totalSamples := streamer.Len()
	if totalSamples == 0 {
		return nil, fmt.Errorf("no audio data found in file: %s", path)
	}

	samplesPerPoint := totalSamples / 1000
	if samplesPerPoint == 0 {
		samplesPerPoint = 1
	}

	waveform := make([][2]float64, 0, 1000)
	buffer := make([][2]float64, samplesPerPoint)

	for i := 0; i < 1000; i++ {
		n, ok := streamer.Stream(buffer)
		if !ok {
			break
		}

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
