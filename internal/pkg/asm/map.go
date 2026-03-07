package asm

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/laghoule/tiled2map/internal/pkg/tiled"
)

type mapHeader struct {
	width        uint8
	height       uint8
}

const (
	mapExt        = ".map"
	mapHeaderSize = 2 // width, height
)

func (a *ASMLinker) createMap(sceneDimension Dimension) error {
	bg, err := a.TileMap.GetLayer(tiled.BackgroundLayerName)
	if err != nil {
		return fmt.Errorf("failed to get background layer: %w", err)
	}

	fg, err := a.TileMap.GetLayer(tiled.ForegroundLayerName)
	if err != nil {
		return fmt.Errorf("failed to get foreground layer: %w", err)
	}

	bgMap := a.convertToMap(sceneDimension, bg)
	fgMap := a.convertToMap(sceneDimension, fg)
	mergedMap := slices.Concat(bgMap, fgMap)

	if err = a.writeMap(sceneDimension, mergedMap); err != nil {
		return err
	}

	return nil
}

// convertToMap converts a tiled layer to a flat localID uint8 map
func (a *ASMLinker) convertToMap(sceneDimension Dimension, layer *tiled.Layer) []uint8 {
	m := make([]uint8, 0, layer.Width*layer.Height)
	for y := 0; y < layer.Height; y += sceneDimension.Height {
		for x := 0; x < layer.Width; x += sceneDimension.Width {

			// Iterate over the scene dimensions
			for sceneY := y; sceneY < y+sceneDimension.Height; sceneY++ {
				for sceneX := x; sceneX < x+sceneDimension.Width; sceneX++ {
					// (y * width) + x
					gid := layer.Data[sceneX+sceneY*layer.Width]
					localID := a.GIDToLocalID[gid]
					m = append(m, localID)
				}
			}

		}
	}

	return m
}

// writeMap writes a flat localID uint8 map to a file.
func (a *ASMLinker) writeMap(sceneDimension Dimension, m []uint8) error {
	fileName := filepath.Join(a.FileOutput.Path, fmt.Sprintf("%s-wrld%s", a.FileOutput.FilePrefix, mapExt))

	mapFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create map file: %w", err)
	}
	defer mapFile.Close()

	// Write the map header
	header := mapHeader{
		width:        uint8(sceneDimension.Width),
		height:       uint8(sceneDimension.Height),
	}

	if err := writeMapHeader(mapFile, header); err != nil {
		return fmt.Errorf("failed to write map header: %w", err)
	}

	for data := range m {
		err := binary.Write(mapFile, binary.LittleEndian, m[data])
		if err != nil {
			return fmt.Errorf("failed to write map data: %w", err)
		}
	}

	return nil
}

// writeMapHeader writes a map header to a file.
func writeMapHeader(file *os.File, header mapHeader) error {
	if err := binary.Write(file, binary.LittleEndian, header); err != nil {
		return fmt.Errorf("failed to write map header to file %s: %w", file.Name(), err)
	}
	return nil
}
