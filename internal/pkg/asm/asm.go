package asm

import (
	"github.com/laghoule/tiled2map/internal/pkg/tiled"
)

// CreateAndSave generates the ASM tile references file
func CreateAndSave(filePrefix string, tilesInfo []tiled.TileInfo) error {
	if err := createTilesRefs(filePrefix, tilesInfo); err != nil {
		return err
	}

	return nil
}
