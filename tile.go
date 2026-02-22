package main

// TileInfo represents the information about a tile, including its source image and position in the tileset
type TileInfo struct {
	SourceImage string
	GID         int
	LocalID     int
	X, Y        int
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
