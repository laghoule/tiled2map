package atlas

import (
	"fmt"
	"image"
	"image/color"
	"os"
)

// getPaletteFromPNG loads the image from the specified source and returns its palette
func getPaletteFromPNG(imgPath string) (*image.Paletted, error) {
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %v", err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image file: %v", err)
	}

	return img.(*image.Paletted), nil
}

// arePaletteEqual compares two palettes and returns true if they are equal, false otherwise
func arePaletteEqual(p1, p2 color.Palette) bool {
	if len(p1) != len(p2) {
		return false
	}

	for i := range p1 {
		r1, g1, b1, _ := p1[i].RGBA()
		r2, g2, b2, _ := p2[i].RGBA()

		if r1 != r2 || g1 != g2 || b1 != b2 {
			return false
		}
	}

	return true
}
