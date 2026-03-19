package main

import (
	"Easy-Tailor/pkg/waveform"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/wav"
)

// App struct
type App struct {
	ctx context.Context
}

type FileInfo struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Path       string `json:"path"`
	DurationMs int64  `json:"durationMs"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetFileInfo(filePath string) (*FileInfo, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return &FileInfo{}, err
	}

	durationMs, err := getAudioDurationMs(filePath)
	if err != nil {
		return &FileInfo{}, err
	}

	fmt.Println("File info:", info.Name(), info.Size())
	return &FileInfo{
		Name:       info.Name(),
		Size:       info.Size(),
		Path:       filePath,
		DurationMs: durationMs,
	}, nil
}

func (a *App) AnalyzeFile(filePath string) ([][2]float64, error) {
	return waveform.AnalyzeFile(filePath)
}

func getAudioDurationMs(filePath string) (int64, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var (
		streamer beep.StreamSeekCloser
		format   beep.Format
	)

	switch filepath.Ext(filePath) {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".flac":
		streamer, format, err = flac.Decode(f)
	default:
		return 0, fmt.Errorf("unsupported file format: %s", filepath.Ext(filePath))
	}
	if err != nil {
		return 0, err
	}
	defer streamer.Close()

	if streamer.Len() <= 0 {
		return 0, nil
	}

	return format.SampleRate.D(streamer.Len()).Milliseconds(), nil
}
