package asm

import (
	"fmt"
	"os"

	"github.com/laghoule/tiled2map/internal/pkg/tiled"
)

const (
	tileAttribute = "attr"
)

// createTilesRefs generates the ASM tile references file
func createTilesRefs(filePrefix string, tilesInfo []tiled.TileInfo) error {
	filename := fmt.Sprintf("%s.inc", filePrefix)
	asmFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer asmFile.Close()

	index := 0
	fmt.Fprintf(asmFile, "TILES_PROPS LABEL BYTE\n")

	for _, tileInfo := range tilesInfo {
		for _, tile := range tileInfo.Tiles {
			attr := 0.0
			for _, prop := range tile.Properties {
				if prop.Name == tileAttribute {
					ok := false
					attr, ok = prop.Value.(float64)
					if !ok {
						return fmt.Errorf("Invalid attribute value for tile %d: %v\n", tileInfo.GID, prop.Value)
					}
					break
				}
			}
			if attr != 0 {
				fmt.Fprintf(asmFile, " DB %08bb ; Index: %d (GID: %d Source: %s)\n", int(attr), index, tileInfo.GID, tileInfo.SourceImage)
				index++
			}
		}
	}

	return nil
}
