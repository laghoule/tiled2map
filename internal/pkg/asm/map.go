package asm

import (
	"fmt"

	"github.com/laghoule/tiled2map/internal/pkg/tiled"
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

	fmt.Println("bg")
	for _, gid := range bg.Data {
		localID := gidToLocalID[gid]
		fmt.Printf("%d, ", localID)
	}

	fmt.Println()

	fmt.Println("fg")
	for _, gid := range fg.Data {
		localID := gidToLocalID[gid]
		fmt.Printf("%d, ", localID)
	}

	return nil
}
