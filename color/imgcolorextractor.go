package color

import (
	"fmt"
	"image"
	_ "image/jpeg" // need for jpg processing
	"io"

	"github.com/deepakjacob/restyle/logger"
	"github.com/generaltso/vibrant"
	"go.uber.org/zap"
)

// ColoExtractor extracts color from the provided reader
type ColoExtractor interface {
	Colors(r io.Reader) []vibrant.Color
}

// Colors extracts color from the provided reader
func Colors(f io.Reader) error {
	img, _, err := image.Decode(f)
	if err != nil {
		logger.Log.Error("image decoding failed", zap.Error(err))
		return err
	}

	// palette, err := vibrant.NewPaletteFromImage(img)
	palette, err := vibrant.NewPalette(img, 10)
	if err != nil {
		logger.Log.Error("error in creating palette", zap.Error(err))
		return err
	}

	for name, swatch := range palette.ExtractAwesome() {
		fmt.Printf("/* %s (population: %d) */\n%s\n\n", name, swatch.Population, swatch)
	}
	return nil
}
