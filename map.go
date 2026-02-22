package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// REF: https://doc.mapeditor.org/en/stable/reference/json-map-format/#

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

// NewMap reads a Tiled map file in JSON format and unmarshals it into a Map struct
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
