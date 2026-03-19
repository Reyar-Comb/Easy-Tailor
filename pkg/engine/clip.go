package engine

import (
	"github.com/gopxl/beep/v2"
)

type clip struct {
	ID     string
	Name   string
	Buffer *beep.Buffer
	Format beep.Format

	Edit Edit
}

type Edit struct {
	TrimStart float64
	TrimEnd   float64
}
