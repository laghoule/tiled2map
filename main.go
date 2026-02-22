package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// Map represents the structure of a Tiled map
type Map struct {
	Type       string    `json:"type"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	TileHeight int       `json:"tileheight"`
	TileWidth  int       `json:"tilewidth"`
	Layers     []Layer   `json:"layers"`
	TileSets   []TileSet `json:"tilesets"`
}

// Layer represents a layer in the Tiled map
type Layer struct {
	Class  string `json:"class"`
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Data   []int  `json:"data"`
}

// TileSet represents a tileset in the Tiled map
type TileSet struct {
	FirstGID         int    `json:"firstgid"`
	Name             string `json:"name"`
	Image            string `json:"image"`
	Columns          int    `json:"columns"`
	TileHeight       int    `json:"tileheight"`
	TileWidth        int    `json:"tilewidth"`
	TileCount        int    `json:"tilecount"`
	TransparentColor string `json:"transparentcolor"`
}

// TileInfo represents the information about a tile, including its source image and position in the tileset
type TileInfo struct {
	SourceImage string
	GID         int
	LocalID     int
	X, Y        int
}

var (
	version   = "unknown"
	gitCommit = "unknown"
)

func NewMap(mapFile string) (*Map, error) {
	mapData, err := os.ReadFile(mapFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read map file: %v", err)
	}

	m := &Map{}
	if err := json.Unmarshal(mapData, m); err != nil {
		return nil, fmt.Errorf("failed to parse map file: %v", err)
	}

	return m, nil
}

// gettUniqueGID retrieves all unique GIDs used in the map layers, excluding GID 0 which represents empty tiles
func getUniqueGID(layers []Layer) []int {
	uniqueGID := make(map[int]bool)
	for _, layer := range layers {
		for _, gid := range layer.Data {
			if gid != 0 {
				uniqueGID[gid] = true
			}
		}
	}

	gids := []int{}
	for gid, _ := range uniqueGID {
		gids = append(gids, gid)
	}

	return gids
}

// findTileSet finds the appropriate tileset for a given GID
func findTileSet(gid int, tileSet []TileSet) *TileSet {
	best := &TileSet{}

	for _, ts := range tileSet {
		if gid >= ts.FirstGID && gid <= (ts.FirstGID+ts.TileCount) {
			best = &ts
			break
		}
	}

	return best
}

// getTilesInfo retrieves the tile information for each unique GID, including the source image and position in the tileset
func getTilesInfo(allGIDs []int, tilesSet []TileSet) []TileInfo {
	tilesInfo := []TileInfo{}
	for _, gid := range allGIDs {
		ts := findTileSet(gid, tilesSet)

		if ts != nil {
			tilesInfo = append(tilesInfo, TileInfo{
				SourceImage: ts.Image,
				GID:         gid,
				LocalID:     gid - ts.FirstGID,
				X:           (gid - ts.FirstGID) % ts.Columns * ts.TileWidth,
				Y:           (gid - ts.FirstGID) / ts.Columns * ts.TileHeight,
			})
		}
	}

	return tilesInfo
}

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

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
