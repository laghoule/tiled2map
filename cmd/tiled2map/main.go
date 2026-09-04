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
	fmt.Printf("Tiled2map version: %s, git commit: %s\n", version, gitCommit)
	mapFile := flag.String("map", "", "Path to the Tiled map file (JSON format)")
	destPath := flag.String("dest", ".", "Destination path for the generated files")
	gridDimension := flag.String("dimension", "", "Dimension of the grid in the map (width x height)")
	generateAttributes := flag.Bool("attributes", true, "Generate attributes for the tiles") // TODO: enable when scene mode detected?
	filePrefix := flag.String("fileprefix", "master", "Prefix for the generated files")
	flag.Parse()

	if *mapFile == "" {
		flag.Usage()
		exitWithError(fmt.Errorf("map file not specified"))
	}

	if err := validateDestPath(*destPath); err != nil {
		exitWithError(err)
	}

	if *gridDimension == "" {
		exitWithError(fmt.Errorf("grid dimension not specified"))
	}

	tileMap, err := tiled.NewMap(*mapFile)
	if err != nil {
		exitWithError(err)
	}

	allGIDs := tiled.GetUniqueGID(tileMap.Layers)
	tilesInfo := tiled.GetSortedTilesInfo(allGIDs, tileMap.TileSets)
	gidLocalTIL := tiled.GetGIDToLocalTIL(allGIDs)

	gridDim, err := asm.ExtractDimension(*gridDimension)
	if err != nil {
		exitWithError(err)
	}

	if len(tilesInfo) == 0 {
		exitWithError(fmt.Errorf("no tiles found"))
	}

	numScene := (tileMap.Width * tileMap.Height) / (gridDim.Width * gridDim.Height)
	if numScene > 1 {
		fmt.Println("Scene mode enabled")
		fmt.Printf("\nNumber of scenes: %d\n", numScene)
		fmt.Printf("Scene dimension: %dx%d\n", gridDim.Width, gridDim.Height)
		fmt.Printf("Scene size: %d bytes\n", gridDim.Width*gridDim.Height)
	}

	numTiles := len(tilesInfo)
	tileSize := tilesInfo[0].Dimension.Width * tilesInfo[0].Dimension.Height
	tilesetSize := numTiles * tileSize
	fmt.Printf("\nNumber of tiles: %d\n", numTiles)
	fmt.Printf("Tile dimension: %dx%d\n", tilesInfo[0].Dimension.Width, tilesInfo[0].Dimension.Height)
	fmt.Printf("Tile size: %d bytes\n", tileSize)
	fmt.Printf("Tileset size: %d bytes\n", tilesetSize)

	master, err := atlas.NewMaster(*destPath, *filePrefix, tilesInfo)
	if err != nil {
		exitWithError(err)
	}

	// Create and save the master atlas file
	err = master.CreateAndSave()
	if err != nil {
		exitWithError(err)
	}

	// Create and save the ASM file with the extracted scene dimension
	asmLinker := asm.NewASMLinker(*destPath, *filePrefix, tileMap, tilesInfo, gidLocalTIL, *generateAttributes)
	err = asmLinker.CreateAndSave(gridDim)
	if err != nil {
		exitWithError(err)
	}

	fmt.Println("\nDone!")
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
