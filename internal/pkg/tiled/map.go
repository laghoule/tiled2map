package tiled

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	BackgroundLayerName = "bg"
	ForegroundLayerName = "fg"
	BoundLayerName      = "bound"
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
	Tiles            []Tile `json:"tiles"`
	TileHeight       int    `json:"tileheight"`
	TileWidth        int    `json:"tilewidth"`
	TileCount        int    `json:"tilecount"`
	TransparentColor string `json:"transparentcolor"`
}

// Tile represents a tile in the tileset, including its properties
type Tile struct {
	ID         int        `json:"id"`
	Properties []Property `json:"properties"`
}

// Property represents a custom property of a tile, which can be of any type (string, int, bool, etc.)
type Property struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value any    `json:"value"`
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

	if err := m.validate(); err != nil {
		return nil, fmt.Errorf("map validation failed: %v", err)
	}

	return m, nil
}

// getProperties retrieves the properties of a tile based on its localID by finding the corresponding tileset and tile information
func getTileProperties(localID int, tileset *TileSet) []Tile {
	tiles := []Tile{}

	for _, tile := range tileset.Tiles {
		if tile.ID != localID {
			continue
		}

		if len(tile.Properties) > 0 {
			tiles = append(tiles, Tile{
				ID:         tile.ID,
				Properties: tile.Properties,
			})
		}
	}

	return tiles
}

// GettUniqueGID retrieves all unique GIDs used in the map layers, excluding GID 0 which represents empty tiles
func GetUniqueGID(layers []Layer) []int {
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

// validate checks if the map contains the required layers (background and foreground) and that their names are valid
func (m *Map) validate() error {
	if len(m.Layers) == 0 {
		return fmt.Errorf("map must contain 2 or 3 layer: %s, %s, %s (optional)", BackgroundLayerName, ForegroundLayerName, BoundLayerName)
	}

	for _, layer := range m.Layers {
		name := strings.ToLower(layer.Name)
		if name == BoundLayerName {
			continue
		}
		if name != BackgroundLayerName && name != ForegroundLayerName {
			return fmt.Errorf("invalid layer name: %s, expected '%s' or '%s'", layer.Name, BackgroundLayerName, ForegroundLayerName)
		}
	}

	return nil
}

// GetLayer retrieves a layer by its name
func (m *Map) GetLayer(name string) (*Layer, error) {
	for _, layer := range m.Layers {
		if strings.ToLower(layer.Name) == name {
			return &layer, nil
		}
	}

	return nil, fmt.Errorf("layer not found: %s", name)
}
