package asm

import (
	"fmt"

	"github.com/laghoule/tiled2map/internal/pkg/tiled"
)

const (
	includeExt = ".inc"
)

// ASMLinker links the map to scene and language generated references
type ASMLinker struct {
	FilePrefix   string
	TileMap      *tiled.Map
	TilesInfo    []tiled.TileInfo
	GIDToLocalID tiled.GIDToLocalTIL
}

// Dimension represents the dimensions of a map.
type Dimension struct {
	Width  int
	Height int
}

// NewASMLinker creates a new ASMLinker.
func NewASMLinker(filePrefix string, tileMap *tiled.Map, tilesInfo []tiled.TileInfo, gidToLocalID tiled.GIDToLocalTIL) *ASMLinker {
	return &ASMLinker{
		FilePrefix:   filePrefix,
		TileMap:      tileMap,
		TilesInfo:    tilesInfo,
		GIDToLocalID: gidToLocalID,
	}
}

// CreateAndSave creates the assembly files, the map and saves them to disk.
func (a *ASMLinker) CreateAndSave(sceneDimension Dimension) error {
	if err := a.createTilesRefs(); err != nil {
		return err
	}

	if err := a.createScene(sceneDimension); err != nil {
		return err
	}

	if err := a.createMap(sceneDimension); err != nil {
		return err
	}

	return nil
}

// ExtractDimension extracts the dimension from a string.
func ExtractDimension(dimension string) (Dimension, error) {
	var d Dimension

	_, err := fmt.Sscanf(dimension, "%dx%d", &d.Width, &d.Height)
	if err != nil {
		return Dimension{}, fmt.Errorf("invalid dimension: %s", dimension)
	}

	return d, nil
}
