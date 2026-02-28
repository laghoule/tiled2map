package tiled

import (
	"sort"
)

// TileInfo represents the information about a tile, including its source image and position in the tileset,
// as well as any custom properties defined in the Tiled map
type TileInfo struct {
	SourceImage string
	Dimension   Dimension
	GID         int
	X, Y        int
	Tiles       []Tile
}

// GIDToLocalID represents a mapping between global tile IDs and local tile IDs within a tileset
type GIDToLocalID map[int]int

// Dimension represents the width and height of a tile
type Dimension struct {
	Width  int
	Height int
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

// GetSortedTilesInfo returns a sorted slice of TileInfo objects for the given GIDs and tilesets
func GetSortedTilesInfo(allGIDs []int, tilesSet []TileSet) []TileInfo {
	tilesInfo := []TileInfo{}

	for _, gid := range allGIDs {
		if gid == 0 {
			continue
		}

		ts := findTileSet(gid, tilesSet)

		if ts != nil {
			localID := gid - ts.FirstGID
			tiles := getTileProperties(localID, ts)
			tilesInfo = append(tilesInfo, TileInfo{
				SourceImage: ts.Image,
				GID:         gid,
				Dimension:   Dimension{Width: ts.TileWidth, Height: ts.TileHeight},
				X:           localID % ts.Columns * ts.TileWidth,
				Y:           localID / ts.Columns * ts.TileHeight,
				Tiles:       tiles,
			})

		}
	}

	sort.Slice(tilesInfo, func(i, j int) bool {
		return tilesInfo[i].GID < tilesInfo[j].GID
	})

	return tilesInfo
}

// GetGIDToLocalID returns a map of global IDs to local IDs for the given GIDs and tilesets
func GetGIDToLocalID(allGIDs []int, tilesSet []TileSet) GIDToLocalID {
	g2l := make(GIDToLocalID)

	for _, gid := range allGIDs {
		ts := findTileSet(gid, tilesSet)

		if ts != nil {
			localID := gid - ts.FirstGID
			g2l[gid] = localID
		}
	}

	return g2l
}
