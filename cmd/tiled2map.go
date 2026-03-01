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
	assetsName := flag.String("name", "master", "Name of the assets to output")
	debug := flag.Bool("debug", false, "Enable debug mode to print additional information")
	flag.Parse()

	m, err := tiled.NewMap(*mapFile)
	if err != nil {
		exitWithError(err)
	}

	allGIDs := tiled.GetUniqueGID(m.Layers)
	tilesInfo := tiled.GetSortedTilesInfo(allGIDs, m.TileSets)
	gidLocalID := tiled.GetGIDToLocalID(allGIDs, m.TileSets)

	if *debug {
		for _, tileInfo := range tilesInfo {
			fmt.Printf("Tile GID: %d\n Source Image: %s, Local ID: %d, X: %d, Y: %d\n Tiles: %v\n",
				tileInfo.GID, tileInfo.SourceImage, gidLocalID[tileInfo.GID], tileInfo.X, tileInfo.Y, tileInfo.Tiles)
			fmt.Println()
		}
	}

	master, err := atlas.NewMaster(tilesInfo)
	if err != nil {
		exitWithError(err)
	}

	err = master.CreateAndSave(*assetsName)
	if err != nil {
		exitWithError(err)
	}

	err = asm.CreateAndSave(m, *assetsName, tilesInfo)
	if err != nil {
		exitWithError(err)
	}

	err = asm.CreateMap(m, gidLocalID)
	if err != nil {
		exitWithError(err)
	}	
}

// exitWithError prints the error message to standard error and exits the program with a non-zero status code
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
