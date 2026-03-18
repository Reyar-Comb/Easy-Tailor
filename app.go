package main

import (
	"Easy-Tailor/pkg/waveform"
	"context"
	"fmt"
	"os"
)

// App struct
type App struct {
	ctx context.Context
}

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Path string `json:"path"`
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
	fmt.Println("File info:", info.Name(), info.Size())
	return &FileInfo{
		Name: info.Name(),
		Size: info.Size(),
		Path: filePath,
	}, nil
}

func (a *App) AnalyzeFile(filePath string) ([][2]float64, error) {
	return waveform.AnalyzeFile(filePath)
}
