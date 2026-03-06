package asm

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// SceneTemplateData represents the data required to generate a single scene entry.
type SceneTemplateData struct {
	Name      string
	MapOffset int
	NorthName string
	SouthName string
	EastName  string
	WestName  string
	MusicName string
}

// SceneTemplatePayload is the top-level data passed to the scene template.
type SceneTemplatePayload struct {
	Prefix     string
	BufferSize int
	Scenes     []SceneTemplateData
}

// createScene generates a scene template based on the provided dimension.
func (a *ASMLinker) createScene(sceneDimension Dimension) error {
	tpl, err := template.ParseFiles("tmpl/scene.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse scene template: %w", err)
	}

	scene := []SceneTemplateData{}
	sceneTiles := sceneDimension.Width * sceneDimension.Height
	numScenesX := int(a.TileMap.Width) / sceneDimension.Width
	numScenesY := int(a.TileMap.Height) / sceneDimension.Height

	// Scene neighbor helper
	getNeighbor := func(sx, sy int, cond bool) string {
		if cond {
			return fmt.Sprintf("OFFSET %s_scene_%d_%d", a.FileOutput.FilePrefix, sx, sy)
		}
		return "0"
	}

	for y := range int(a.TileMap.Height) / sceneDimension.Height {
		for x := range int(a.TileMap.Width) / sceneDimension.Width {
			// offset is the 2D -> 1D transformation
			currentOffset := ((y*numScenesX)+x)*sceneTiles + int(mapHeaderSize)

			scene = append(scene, SceneTemplateData{
				Name:      fmt.Sprintf("%s_scene_%d_%d", a.FileOutput.FilePrefix, x, y),
				MapOffset: currentOffset,
				NorthName: getNeighbor(x, y-1, y > 0),
				SouthName: getNeighbor(x, y+1, y < numScenesY-1),
				EastName:  getNeighbor(x+1, y, x < numScenesX-1),
				WestName:  getNeighbor(x-1, y, x > 0),
				MusicName: fmt.Sprintf("%s_music_%d_%d", a.FileOutput.FilePrefix, x, y),
			})
		}
	}

	filename := filepath.Join(a.FileOutput.Path, fmt.Sprintf("%s-scne%s", a.FileOutput.FilePrefix, includeExt))
	sceneFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create scene file: %w", err)
	}
	defer sceneFile.Close()

	// BufferSize is the total number of bytes needed for bg/fg map buffers:
	// all tiles across the entire map, plus the 2-byte header per file.
	bufferSize := a.TileMap.Width*a.TileMap.Height + mapHeaderSize

	payload := SceneTemplatePayload{
		Prefix:     a.FileOutput.FilePrefix,
		BufferSize: bufferSize,
		Scenes:     scene,
	}

	err = tpl.Execute(sceneFile, payload)
	if err != nil {
		return fmt.Errorf("failed to execute scene template: %w", err)
	}

	return nil
}
