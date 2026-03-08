package atlas

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

// Create generates the master image by drawing each tile onto it.
func (m *Master) createIMG() error {
	imageCache := make(map[string]image.Image)

	for i, tile := range m.Tiles {
		src, err := getOrLoadImage(tile.SourceImage, imageCache)
		if err != nil {
			return fmt.Errorf("failed to load image: %v", err)
		}

		tilePalleted, ok := src.(*image.Paletted)
		if !ok {
			return fmt.Errorf("tile %d, from %s is not a paletted image", i, tile.SourceImage)
		}

		if !arePaletteEqual(m.Palette, tilePalleted.Palette) {
			return fmt.Errorf("each tile must have the same palette. Tile %d, from %s has a different palette", i, tile.SourceImage)
		}

		// extract tile from the image
		tileRect := image.Rect(tile.X, tile.Y, tile.X+m.Dimension.Width, tile.Y+m.Dimension.Height)

		// calculate the destination rectangle for the tile in the master image
		destRect := image.Rect(0, i*m.Dimension.Height, m.Dimension.Width, (i+1)*m.Dimension.Height)

		// draw the tile onto the master image
		draw.Draw(m.Image, destRect, src, tileRect.Min, draw.Src)
	}

	// Create the raw image
	for y := range m.Dimension.Height * m.TileCount {
		for x := range m.Dimension.Width {
			m.RawImage = append(m.RawImage, m.Image.ColorIndexAt(x, y))
		}
	}

	return nil
}

// SavePNG saves the master image to a file in PNG format.
func (m *Master) savePNG() error {
	filename := filepath.Join(m.FileOutput.Path, fmt.Sprintf("%s-ts.png", m.FileOutput.FilePrefix))
	masterFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create master file: %v", err)
	}
	defer masterFile.Close()

	if err := png.Encode(masterFile, m.Image); err != nil {
		return fmt.Errorf("failed to encode master image: %v", err)
	}

	return nil
}

// SaveTIL saves the raw image to a file in TIL format.
func (m *Master) saveTIL() error {
	filename := filepath.Join(m.FileOutput.Path, fmt.Sprintf("%s-ts.til", m.FileOutput.FilePrefix))
	tilFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create raw file: %v", err)
	}
	defer tilFile.Close()

	if _, err = tilFile.Write(m.RawImage); err != nil {
		return fmt.Errorf("failed to write data to file %s: %v", filename, err)
	}

	return nil
}
