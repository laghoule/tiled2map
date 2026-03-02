package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/laghoule/tiled2map/internal/pkg/asm"
	"github.com/laghoule/tiled2map/internal/pkg/atlas"
	"github.com/laghoule/tiled2map/internal/pkg/tiled"
)

var (
	version   = "unknown"
	gitCommit = "unknown"
)

func main() {
	fmt.Printf("tile2map version: %s, git commit: %s\n", version, gitCommit)
	mapFile := flag.String("map", "", "Path to the Tiled map file (JSON format)")
	sceneDimension := flag.String("dimension", "20x11", "Dimension of each scenes")
	filePrefix := flag.String("fileprefix", "master", "Prefix for the generated files")
	flag.Parse()

	tileMap, err := tiled.NewMap(*mapFile)
	if err != nil {
		exitWithError(err)
	}

	allGIDs := tiled.GetUniqueGID(tileMap.Layers)
	tilesInfo := tiled.GetSortedTilesInfo(allGIDs, tileMap.TileSets)
	gidLocalID := tiled.GetGIDToLocalID(allGIDs, tileMap.TileSets)

	master, err := atlas.NewMaster(tilesInfo)
	if err != nil {
		exitWithError(err)
	}

	// Create and save the master atlas file
	err = master.CreateAndSave(*filePrefix)
	if err != nil {
		exitWithError(err)
	}

	// Extract scene dimension from the command line argument
	dimension, err := asm.ExtractDimension(*sceneDimension)
	if err != nil {
		exitWithError(err)
	}

	// Create and save the ASM file with the extracted scene dimension
	asmLinker := asm.NewASMLinker(*filePrefix, tileMap, tilesInfo, gidLocalID)
	err = asmLinker.CreateAndSave(dimension)
	if err != nil {
		exitWithError(err)
	}
}

// exitWithError prints the error message to standard error and exits the program with a non-zero status code
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
