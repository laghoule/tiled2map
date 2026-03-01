package asm

import (
	"github.com/laghoule/tiled2map/internal/pkg/tiled"
)

const (
	includeExt = ".inc"
)

// dimension represents the dimensions of a map.
type dimension struct {
	width  uint8
	height uint8
}

// CreateAndSave generates the ASM tile references file
func CreateAndSave(mp *tiled.Map, filePrefix string, tilesInfo []tiled.TileInfo) error {
	if err := createTilesRefs(filePrefix, tilesInfo); err != nil {
		return err
	}

	d := dimension{
		width:  uint8(mp.Width),
		height: uint8(mp.Height),
	}

	if err := createScene(filePrefix, d); err != nil {
		return err
	}

	return nil
}
