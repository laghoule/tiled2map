package atlas

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"os"
	"sort"

	"github.com/laghoule/tiled2map/internal/pkg/tiled"
)

const (
	tileSpacing = 1
)

// Master represents the master image that contains all the tiles from the Tiled map. It includes the image itself and the raw byte data of the image.
type Master struct {
	Tiles     []tiled.TileInfo
	Image     *image.Paletted
	Palette   color.Palette
	TileCount int
	Dimension Dimension
}

// Dimension represents the width and height of a tile
type Dimension struct {
	Width  int
	Height int
}

// TODO: decouple
func NewMaster(tiles []tiled.TileInfo) (*Master, error) {
	width, height := getTileDimension(tiles[0])

	tilesCount := len(tiles)
	masterRect := image.Rect(0, 0, width, height*tilesCount)
	masterImg := image.NewPaletted(masterRect, nil)

	// Get the palette from the first png image
	firstTilePal, err := getPaletteFromPNG(tiles[0].SourceImage)
	if err != nil {
		return nil, fmt.Errorf("failed to get palette: %v", err)
	}

	// Set the palette of the master image to match the first tile's palette
	masterImg.Palette = firstTilePal.Palette

	return &Master{
		Tiles:     tiles,
		Image:     masterImg,
		TileCount: tilesCount,
		Palette:   firstTilePal.Palette,
		Dimension: Dimension{
			Width:  width,
			Height: height,
		},
	}, nil
}

// CreateAndSave creates the master image and saves it to a file in PNG format.
func (m *Master) CreateAndSave(filename string) error {
	if err := m.create(); err != nil {
		return err
	}

	if err := m.save(filename); err != nil {
		return err
	}

	return nil
}

// Create generates the master image by drawing each tile onto it.
func (m *Master) create() error {
	m.Image.Palette = m.Palette

	// Ensure that all tiles are sorted by GID to maintain a consistent order in the master image
	sort.Slice(m.Tiles, func(i, j int) bool {
		return m.Tiles[i].GID < m.Tiles[j].GID
	})

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
		destRect := image.Rect(0, i*m.Dimension.Height, m.Dimension.Width, (i+tileSpacing)*m.Dimension.Height)

		// draw the tile onto the master image
		draw.Draw(m.Image, destRect, src, tileRect.Min, draw.Src)
	}

	return nil
}

// getTileDimension returns the width and height of a tile based on its dimension information
func getTileDimension(tile tiled.TileInfo) (int, int) {
	return tile.Dimension.Width, tile.Dimension.Height
}

// getOrLoadImage retrieves the image from the cache if it exists, or loads it from the file system if it doesn't
func getOrLoadImage(imgPath string, imageCache map[string]image.Image) (image.Image, error) {
	if img, exists := imageCache[imgPath]; exists {
		return img, nil
	}

	imgFile, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %v", err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image file: %v", err)
	}

	imageCache[imgPath] = img

	return img, nil
}

// Save saves the master image to a file in PNG format.
func (m *Master) save(filename string) error {
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
