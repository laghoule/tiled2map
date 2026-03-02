package asm

import (
	"github.com/laghoule/tiled2map/internal/pkg/tiled"
)

const (
	includeExt = ".inc"
)

// dimension represents the dimensions of a map.
type dimension struct {
	width  int
	height int
}

// CreateAndSave generates the ASM tile references file
func CreateAndSave(m *tiled.Map, filePrefix string, tilesInfo []tiled.TileInfo, gidLocalID tiled.GIDToLocalID) error {
	if err := createTilesRefs(filePrefix, tilesInfo); err != nil {
		return err
	}

	if err := createScene(m, filePrefix); err != nil {
		return err
	}
	
	if err := createMap(m, gidLocalID);  err != nil {
		return err
	}	

	return nil
}
