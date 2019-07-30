package color

import (
	"fmt"
	"image"
	"io"

	"github.com/generaltso/vibrant"
)

var (
	errColorExtraction = fmt.Errorf("error extracting colors from image")
)

// Color extracts color from the provided image
func Color(f io.Reader) error {
	img, _, err := image.Decode(f)
	if err != nil {
		return errColorExtraction
	}

	palette, err := vibrant.NewPaletteFromImage(img)
	if err != nil {
		return errColorExtraction
	}

	for name, swatch := range palette.ExtractAwesome() {
		fmt.Printf("/* %s (population: %d) */\n%s\n\n", name, swatch.Population, swatch)
	}
	return nil
}
