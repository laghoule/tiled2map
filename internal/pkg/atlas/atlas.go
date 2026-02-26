package atlas

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"

	"github.com/laghoule/tiled2map/internal/pkg/tiled"
)

const (
	tileSpacing = 1
)

// Master represents the master image that contains all the tiles from the Tiled map. It includes the image itself and the raw byte data of the image.
type Master struct {
	Tiles     []tiled.TileInfo
	Image     *image.Paletted
	RawImage  []byte
	Palette   color.Palette
	TileCount int
	Dimension Dimension
}

// Dimension represents the width and height of a tile
type Dimension struct {
	Width  int
	Height int
}

// NewMaster creates a new Master instance with the given tiles.
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
		RawImage:  []byte{},
		TileCount: tilesCount,
		Palette:   firstTilePal.Palette,
		Dimension: Dimension{
			Width:  width,
			Height: height,
		},
	}, nil
}

// CreateAndSave creates the master image and saves it to a file in PNG format.
func (m *Master) CreateAndSave(filePrefix string) error {
	if err := m.createIMG(); err != nil {
		return err
	}

	if err := m.savePNG(filePrefix); err != nil {
		return err
	}
	
	if err := m.saveTIL(filePrefix); err != nil {
		return err
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
