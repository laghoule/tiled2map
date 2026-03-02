package asm

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/laghoule/tiled2map/internal/pkg/tiled"
)

const (
	mapExt        = ".map"
	mapHeaderSize = uint8(2) // width + height
)

func (a *ASMLinker) createMap() error {
	bg, err := a.TileMap.GetLayer(tiled.BackgroundLayerName)
	if err != nil {
		return fmt.Errorf("failed to get background layer: %w", err)
	}

	fg, err := a.TileMap.GetLayer(tiled.ForegroundLayerName)
	if err != nil {
		return fmt.Errorf("failed to get foreground layer: %w", err)
	}

	bgMap := a.convertToMap(bg)
	fgMap := a.convertToMap(fg)

	if err = a.writeMap("bg-world", bgMap); err != nil {
		return fmt.Errorf("failed to write bg map: %w", err)
	}

	if err = a.writeMap("fg-world", fgMap); err != nil {
		return fmt.Errorf("failed to write fg map: %w", err)
	}

	return nil
}

// convertToMap converts a tiled layer to a flat localID uint8 map
func (a *ASMLinker) convertToMap(l *tiled.Layer) []uint8 {
	m := make([]uint8, 0, l.Width*l.Height)
	for y := 0; y < l.Height; y += a.Dimension.Height {
		for x := 0; x < l.Width; x += a.Dimension.Width {

			// Iterate over the scene dimensions
			for sceneY := y; sceneY < y+a.Dimension.Height; sceneY++ {
				for sceneX := x; sceneX < x+a.Dimension.Width; sceneX++ {
					// (y * width) + x
					gid := l.Data[sceneX+sceneY*l.Width]
					localID := a.GIDToLocalID[gid]
					m = append(m, localID)
				}
			}

		}
	}

	return m
}

// writeMap writes a flat localID uint8 map to a file.
func (a *ASMLinker) writeMap(layer string, m []uint8) error {
	fileName := fmt.Sprintf("%s-%s%s", a.FilePrefix, layer, mapExt)

	mapFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create map file: %w", err)
	}
	defer mapFile.Close()

	// Write map dimensions as header to map file
	binary.Write(mapFile, binary.LittleEndian, uint8(a.Dimension.Width))
	binary.Write(mapFile, binary.LittleEndian, uint8(a.Dimension.Height))

	for data := range m {
		err := binary.Write(mapFile, binary.LittleEndian, m[data])
		if err != nil {
			return fmt.Errorf("failed to write map data: %w", err)
		}
	}

	return nil
}
