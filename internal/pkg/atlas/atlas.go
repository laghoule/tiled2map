package atlas

import (
	"fmt"
	"image"
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
	Image    *image.Paletted
	RawImage []byte
}

func NewMaster(tiles []tiled.TileInfo) (*Master, error) {
	width, height := getTileDimension(tiles[0])

	count := len(tiles)
	masterRect := image.Rect(0, 0, width, height*count)
	masterImg := image.NewPaletted(masterRect, nil)

	rawBytes := make([]byte, width*height)
	imageCache := make(map[string]image.Image)

	// Get the palette from the first png image
	firstTilePal, err := getPaletteFromPNG(tiles[0].SourceImage)
	if err != nil {
		return nil, fmt.Errorf("failed to get palette: %v", err)
	}

	// Set the palette of the master image to match the first tile's palette
	masterImg.Palette = firstTilePal.Palette

	// Ensure that all tiles are sorted by GID to maintain a consistent order in the master image
	sort.Slice(tiles, func(i, j int) bool {
		return tiles[i].GID < tiles[j].GID
	})

	for i, tile := range tiles {
		src, err := getOrLoadImage(tile.SourceImage, imageCache)
		if err != nil {
			return nil, fmt.Errorf("failed to load image: %v", err)
		}

		tilePal, ok := src.(*image.Paletted)
		if !ok {
			return nil, fmt.Errorf("source image is not a paletted image")
		}

		if !arePaletteEqual(firstTilePal, tilePal) {
			return nil, fmt.Errorf("each tile must have the same palette. Tile %d has a different palette", i)
		}

		// extract tile from the image
		tileRect := image.Rect(tile.X, tile.Y, tile.X+width, tile.Y+height)

		// calculate the destination rectangle for the tile in the master image
		destRect := image.Rect(0, i*height, width, (i+tileSpacing)*height)

		// draw the tile onto the master image
		draw.Draw(masterImg, destRect, src, tileRect.Min, draw.Src)
	}

	return &Master{
		Image:    masterImg,
		RawImage: rawBytes,
	}, nil
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
func (m *Master) Save(filename string) error {
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
