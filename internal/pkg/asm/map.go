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

func createMap(m *tiled.Map, gidToLocalID tiled.GIDToLocalID) error {
	bg, err := m.GetLayer(tiled.BackgroundLayerName)
	if err != nil {
		return fmt.Errorf("failed to get background layer: %w", err)
	}

	fg, err := m.GetLayer(tiled.ForegroundLayerName)
	if err != nil {
		return fmt.Errorf("failed to get foreground layer: %w", err)
	}

	// Map dimension
	dim := dimension{width: m.Width, height: m.Height}

	bgMap := convertToMap(dim, bg, gidToLocalID)
	fgMap := convertToMap(dim, fg, gidToLocalID)

	if err = writeMap("bg-world", dim, bgMap); err != nil {
		return fmt.Errorf("failed to write bg map: %w", err)
	}

	if err = writeMap("fg-world", dim, fgMap); err != nil {
		return fmt.Errorf("failed to write fg map: %w", err)
	}

	return nil
}

// convertToMap converts a tiled layer to a flat localID uint8 map
func convertToMap(d dimension, l *tiled.Layer, gidToLocalID tiled.GIDToLocalID) []uint8 {
	m := make([]uint8, 0, l.Width*l.Height)
	for y := 0; y < l.Height; y += d.height {
		for x := 0; x < l.Width; x += d.width {

			// Iterate over the scene dimensions
			for sceneY := y; sceneY < y+d.height; sceneY++ {
				for sceneX := x; sceneX < x+d.width; sceneX++ {
					// (y * width) + x
					gid := l.Data[sceneX+sceneY*l.Width]
					localID := gidToLocalID[gid]
					m = append(m, localID)
				}
			}

		}
	}

	return m
}

// writeMap writes a flat localID uint8 map to a file.
func writeMap(filePrefix string, d dimension, m []uint8) error {
	mapFile, err := os.Create(filePrefix + mapExt)
	if err != nil {
		return fmt.Errorf("failed to create map file: %w", err)
	}
	defer mapFile.Close()

	// Write map dimensions as header to map file
	binary.Write(mapFile, binary.LittleEndian, d.width)
	binary.Write(mapFile, binary.LittleEndian, d.height)

	for data := range m {
		err := binary.Write(mapFile, binary.LittleEndian, m[data])
		if err != nil {
			return fmt.Errorf("failed to write map data: %w", err)
		}
	}

	return nil
}
