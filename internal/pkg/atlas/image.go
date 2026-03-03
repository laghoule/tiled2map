package atlas

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

// tilHeader represents the header of a TIL file.
type tilHeader struct {
	Width      uint8
	Height     uint8
	TilesCount uint8
}

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
			return fmt.Errorf("tile %d is not a paletted image", i)
		}

		if !arePaletteEqual(m.Palette, tilePalleted.Palette) {
			return fmt.Errorf("each tile must have the same palette. Tile %d has a different palette", i)
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
func (m *Master) savePNG(filePrefix string) error {
	filename := fmt.Sprintf("%s-tileset.png", filePrefix)
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
func (m *Master) saveTIL(filePrefix string) error {
	filename := fmt.Sprintf("%s-tileset.til", filePrefix)
	tilFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create raw file: %v", err)
	}
	defer tilFile.Close()

	h := tilHeader{
		Width:      uint8(m.Dimension.Width),
		Height:     uint8(m.Dimension.Height),
		TilesCount: uint8(m.TileCount),
	}

	if err := writeTILHeader(tilFile, h); err != nil {
		return err
	}

	if _, err = tilFile.Write(m.RawImage); err != nil {
		return fmt.Errorf("failed to write data to file %s: %v", filename, err)
	}

	return nil
}

// writeHeader writes the header to the raw file.
func writeTILHeader(file *os.File, header tilHeader) error {
	if err := binary.Write(file, binary.LittleEndian, header); err != nil {
		return fmt.Errorf("failed to write header to file %s: %v", file.Name(), err)
	}
	return nil
}
