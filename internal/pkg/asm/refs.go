package asm

import (
	"fmt"
	"os"

	"github.com/laghoule/tiled2map/internal/pkg/tiled"
)

const (
	tileAttribute = "attr"
)

// CreateAndSave generates the ASM tile references file
func CreateAndSave(filePrefix string, tilesInfo []tiled.TileInfo) error {
	if err := createTilesRefs(filePrefix, tilesInfo); err != nil {
		return err
	}

	return nil
}

// createTilesRefs generates the ASM tile references file
func createTilesRefs(filePrefix string, tilesInfo []tiled.TileInfo) error {
	filename := fmt.Sprintf("%s.inc", filePrefix)
	asmFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer asmFile.Close()

	fmt.Fprintf(asmFile, "TILES_PROPS LABEL BYTE\n")
	for _, tileInfo := range tilesInfo {
		for _, tile := range tileInfo.Tiles {
			attr := 0.0
			for _, prop := range tile.Properties {
				if prop.Name == tileAttribute {
					attr = prop.Value.(float64)
					break
				}
			}
			if attr != 0 {
				fmt.Fprintf(asmFile, " DB %08bb ; GID: %d Source: %s\n", int(attr), tileInfo.GID, tileInfo.SourceImage)
			}
		}
	}

	return nil
}
