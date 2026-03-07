package main

import (
	"errors"
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
	fmt.Printf("tiled2map version: %s, git commit: %s\n", version, gitCommit)
	mapFile := flag.String("map", "", "Path to the Tiled map file (JSON format)")
	destPath := flag.String("dest", ".", "Destination path for the generated files")
	sceneDimension := flag.String("dimension", "20x11", "Dimension of each scenes")
	filePrefix := flag.String("fileprefix", "master", "Prefix for the generated files")
	flag.Parse()

	if *mapFile == "" {
		flag.Usage()
		exitWithError(fmt.Errorf("map file not specified"))
	}

	if err := validateDestPath(*destPath); err != nil {
		exitWithError(err)
	}

	tileMap, err := tiled.NewMap(*mapFile)
	if err != nil {
		exitWithError(err)
	}

	allGIDs := tiled.GetUniqueGID(tileMap.Layers)
	tilesInfo := tiled.GetSortedTilesInfo(allGIDs, tileMap.TileSets)
	gidLocalTIL := tiled.GetGIDToLocalTIL(allGIDs)

	sceneDim, err := asm.ExtractDimension(*sceneDimension)
	if err != nil {
		exitWithError(err)
	}

	fmt.Println()
	fmt.Printf("Number of scenes: %d\n", (tileMap.Width*tileMap.Height)/(sceneDim.Width*sceneDim.Height))
	fmt.Printf("Scene dimension: %dx%d\n", sceneDim.Width, sceneDim.Height)
	fmt.Printf("Scenes size: %d bytes\n", sceneDim.Width*sceneDim.Height)
	fmt.Println()
	fmt.Printf("Number of tiles: %d\n", len(tilesInfo))
	fmt.Printf("Tiles dimension: %dx%d\n", tilesInfo[0].Dimension.Width, tilesInfo[0].Dimension.Height)
	fmt.Printf("Tiles size: %d bytes", tilesInfo[0].Dimension.Width*tilesInfo[0].Dimension.Height)
	fmt.Println("")

	master, err := atlas.NewMaster(*destPath, *filePrefix, tilesInfo)
	if err != nil {
		exitWithError(err)
	}

	// Create and save the master atlas file
	err = master.CreateAndSave()
	if err != nil {
		exitWithError(err)
	}

	// Extract scene dimension from the command line argument
	dimension, err := asm.ExtractDimension(*sceneDimension)
	if err != nil {
		exitWithError(err)
	}

	// Create and save the ASM file with the extracted scene dimension
	asmLinker := asm.NewASMLinker(*destPath, *filePrefix, tileMap, tilesInfo, gidLocalTIL)
	err = asmLinker.CreateAndSave(dimension)
	if err != nil {
		exitWithError(err)
	}

	fmt.Println("")
	fmt.Println("Done!")
}

// validateDestPath validates the destination path for the generated files
func validateDestPath(destPath string) error {
	fInfo, err := os.Lstat(destPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("destination path %s does not exist", destPath)
		}
		return fmt.Errorf("failed to validate destination path: %v", err)
	}
	if !fInfo.IsDir() {
		return fmt.Errorf("destination path is not a directory")
	}
	return nil
}

// exitWithError prints the error message to standard error and exits the program with a non-zero status code
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
