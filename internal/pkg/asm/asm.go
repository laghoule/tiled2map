package asm

import (
	"github.com/laghoule/tiled2map/internal/pkg/tiled"
)

const (
	includeExt = ".inc"
)

// dimension represents the dimensions of a map.
// TODO: uint8 ???
type dimension struct {
	width  uint8
	height uint8
}

// CreateAndSave generates the ASM tile references file
func CreateAndSave(m *tiled.Map, filePrefix string, tilesInfo []tiled.TileInfo) error {
	if err := createTilesRefs(filePrefix, tilesInfo); err != nil {
		return err
	}

	if err := createScene(m, filePrefix); err != nil {
		return err
	}

	return nil
}
