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
	generateAttributes bool
	fileOutput         FileOutput
	tileMap            *tiled.Map
	tilesInfo          []tiled.TileInfo
	gIDToLocalID       tiled.GIDToLocalTIL
}

// Dimension represents the dimensions of a map.
type Dimension struct {
	Width  int
	Height int
}

// FileOutput represents a file output for the atlas.
type FileOutput struct {
	Path       string
	FilePrefix string
}

// NewASMLinker creates a new ASMLinker.
func NewASMLinker(destPath, filePrefix string, tileMap *tiled.Map, tilesInfo []tiled.TileInfo, gidToLocalID tiled.GIDToLocalTIL) *ASMLinker {
	return &ASMLinker{
		fileOutput: FileOutput{
			Path:       destPath,
			FilePrefix: filePrefix,
		},
		tileMap:      tileMap,
		tilesInfo:    tilesInfo,
		gIDToLocalID: gidToLocalID,
	}
}

// CreateAndSave creates the assembly files, the map and saves them to disk.
func (a *ASMLinker) CreateAndSave(sceneDimension Dimension) error {
	a.generateAttributes = sceneDimension != (Dimension{})
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
	var dim Dimension

	if dimension == "" {
		return dim, nil
	}

	_, err := fmt.Sscanf(dimension, "%dx%d", &dim.Width, &dim.Height)
	if err != nil {
		return Dimension{}, fmt.Errorf("invalid dimension: %s", dimension)
	}

	return dim, nil
}
