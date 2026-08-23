package asm

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed tmpl/*
var tmplFS embed.FS

const (
	sceneTmplFile = "tmpl/scene.tmpl"
	layersLen     = 2 // bg & fg
)

// SceneData represents the data required to generate a single scene entry.
type SceneData struct {
	Name      string
	MapOffset int
	NorthName string
	SouthName string
	EastName  string
	WestName  string
}

// MapData represents the data required to generate the map data.
type MapData struct {
	Dimension Dimension
	LayerSize int
}

// SceneTemplatePayload is the top-level data passed to the scene template.
type SceneTemplatePayload struct {
	Prefix     string
	BufferSize int

	MapData        MapData
	SceneDimension Dimension
	Scenes         []SceneData
}

// createScene generates a scene template based on the provided dimension.
func (a *ASMLinker) createScene(sceneDimension Dimension) error {
	if sceneDimension == (Dimension{}) {
		return nil
	}

	scenes := []SceneData{}
	sceneTiles := sceneDimension.Width * sceneDimension.Height
	numScenesX := int(a.tileMap.Width) / sceneDimension.Width
	numScenesY := int(a.tileMap.Height) / sceneDimension.Height

	// Scene neighbor helper
	getNeighbor := func(sx, sy int, cond bool) string {
		if cond {
			return fmt.Sprintf("OFFSET %s_scene_%d_%d", a.fileOutput.FilePrefix, sx, sy)
		}
		return "0"
	}

	for y := range int(a.tileMap.Height) / sceneDimension.Height {
		for x := range int(a.tileMap.Width) / sceneDimension.Width {
			// offset is the 2D -> 1D transformation
			currentOffset := ((y * numScenesX) + x) * sceneTiles

			scenes = append(scenes, SceneData{
				Name:      fmt.Sprintf("%s_scene_%d_%d", a.fileOutput.FilePrefix, x, y),
				MapOffset: currentOffset,
				NorthName: getNeighbor(x, y-1, y > 0),
				SouthName: getNeighbor(x, y+1, y < numScenesY-1),
				EastName:  getNeighbor(x+1, y, x < numScenesX-1),
				WestName:  getNeighbor(x-1, y, x > 0),
			})
		}
	}

	filename := filepath.Join(a.fileOutput.Path, fmt.Sprintf("%s-scne%s", a.fileOutput.FilePrefix, includeExt))
	sceneFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create scene file: %w", err)
	}
	defer sceneFile.Close()

	// BufferSize is the total number of bytes needed for bg/fg map buffers
	bufferSize := ((a.tileMap.Width * a.tileMap.Height) * layersLen)
	mapLayerSize := a.tileMap.Width * a.tileMap.Height

	payload := SceneTemplatePayload{
		Prefix:     a.fileOutput.FilePrefix,
		BufferSize: bufferSize,
		MapData: MapData{
			Dimension: Dimension{
				Width:  a.tileMap.Width,
				Height: a.tileMap.Height,
			},
			LayerSize: mapLayerSize,
		},
		SceneDimension: Dimension{
			Width:  sceneDimension.Width,
			Height: sceneDimension.Height,
		},
		Scenes: scenes,
	}

	funcMap := template.FuncMap{
		"toUpper": strings.ToUpper,
	}

	tpl, err := template.New("scene.tmpl").Funcs(funcMap).ParseFS(tmplFS, sceneTmplFile)
	if err != nil {
		return fmt.Errorf("failed to parse scene template: %w", err)
	}

	err = tpl.Execute(sceneFile, payload)
	if err != nil {
		return fmt.Errorf("failed to execute scene template: %w", err)
	}

	return nil
}
