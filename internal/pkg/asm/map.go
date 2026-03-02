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

func CreateMap(m *tiled.Map, gidToLocalID tiled.GIDToLocalID) error {
	bg, err := m.GetLayer(tiled.BackgroundLayerName)
	if err != nil {
		return fmt.Errorf("failed to get background layer: %w", err)
	}

	fg, err := m.GetLayer(tiled.ForegroundLayerName)
	if err != nil {
		return fmt.Errorf("failed to get foreground layer: %w", err)
	}

	d := dimension{
		width:  uint8(m.Width),
		height: uint8(m.Height),
	}

	bgMap := createMap(d, bg, gidToLocalID)
	fgMap := createMap(d, fg, gidToLocalID)

	if err = writeMap("bg-world", d, bgMap); err != nil {
		return fmt.Errorf("failed to write bg map: %w", err)
	}

	if err = writeMap("fg-world", d, fgMap); err != nil {
		return fmt.Errorf("failed to write fg map: %w", err)
	}

	return nil
}

// createMap converts a GID Tiled layer to a flat localID uint8 map.
func createMap(d dimension, l *tiled.Layer, gidToLocalID tiled.GIDToLocalID) []uint8 {
	m := make([]uint8, 0, l.Width*l.Height)
	for y := 0; y < l.Height; y += int(d.height) {
		for x := 0; x < l.Width; x += int(d.width) {

			for sceneY := y; sceneY < y+int(d.height); sceneY++ {
				for sceneX := x; sceneX < x+int(d.width); sceneX++ {
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
