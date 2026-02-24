package tiled

// TileInfo represents the information about a tile, including its source image and position in the tileset,
// as well as any custom properties defined in the Tiled map
// TODO: Implement sort
type TileInfo struct {
	SourceImage string
	GID         int
	LocalID     int
	Dimension   Dimension
	X, Y        int
	Tiles       []Tile
}

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

// GetTilesInfo retrieves the tile information for each unique GID, including the source image and position in the tileset
func GetTilesInfo(allGIDs []int, tilesSet []TileSet) []TileInfo {
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
				LocalID:     localID,
				Dimension:   Dimension{Width: ts.TileWidth, Height: ts.TileHeight},
				X:           localID % ts.Columns * ts.TileWidth,
				Y:           localID / ts.Columns * ts.TileHeight,
				Tiles:       tiles,
			})

		}
	}

	return tilesInfo
}
