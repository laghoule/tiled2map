package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	version   = "unknown"
	gitCommit = "unknown"
)

func main() {
	fmt.Printf("tile2map version: %s, git commit: %s\n", version, gitCommit)
	mapFile := flag.String("map", "", "Path to the Tiled map file (JSON format)")
	flag.Parse()

	m, err := NewMap(*mapFile)
	if err != nil {
		exitWithError(err)
	}

	allGIDs := getUniqueGID(m.Layers)
	tilesInfo := getTilesInfo(allGIDs, m.TileSets)

	for _, tileInfo := range tilesInfo {
		fmt.Printf("Tile GID: %d, Source Image: %s, Local ID: %d, X: %d, Y: %d\n",
			tileInfo.GID, tileInfo.SourceImage, tileInfo.LocalID, tileInfo.X, tileInfo.Y)
	}
}

// exitWithError prints the error message to standard error and exits the program with a non-zero status code
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
