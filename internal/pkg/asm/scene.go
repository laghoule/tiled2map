package asm

import (
	"fmt"
	"os"
	"text/template"
)

// sceneTemplateData represents the data required to generate a scene template.
type SceneTemplateData struct {
	Name      string
	BGOffset  int
	FGOffset  int
	NorthName string
	SouthName string
	EastName  string
	WestName  string
	MusicName string
}

// createScene generates a scene template based on the provided dimension.
func (a *ASMLinker) createScene() error {
	tpl, err := template.ParseFiles("tmpl/scene.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse scene template: %w", err)
	}

	scene := []SceneTemplateData{}
	sceneSize := (int(a.TileMap.Width) / 20) * (int(a.TileMap.Height) / 11)

	// Scene neighbor helper
	getNeighbord := func(sx, sy int, cond bool) string {
		if cond {
			return fmt.Sprintf("OFFSET SCENE_%d_%d", sx, sy)
		}
		return "0"
	}

	for y := range int(a.TileMap.Height) / 11 {
		for x := range int(a.TileMap.Width) / 20 {
			// offset is the 2D -> 1D transformation
			currentOffset := ((y * int(a.TileMap.Width)) + x*sceneSize) + int(mapHeaderSize)

			scene = append(scene, SceneTemplateData{
				Name:      fmt.Sprintf("SCENE_%d_%d", x, y),
				BGOffset:  currentOffset,
				FGOffset:  currentOffset,
				NorthName: getNeighbord(x, y-1, y > 0),
				SouthName: getNeighbord(x, y+1, y < int(a.TileMap.Height)-1),
				EastName:  getNeighbord(x+1, y, x < int(a.TileMap.Width)-1),
				WestName:  getNeighbord(x-1, y, x > 0),
				MusicName: fmt.Sprintf("MUSIC_%d_%d", x, y),
			})
		}
	}

	filename := fmt.Sprintf("%s-scene.%s", a.FilePrefix, includeExt)
	sceneFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create scene file: %w", err)
	}
	defer sceneFile.Close()

	err = tpl.Execute(sceneFile, scene)
	if err != nil {
		return fmt.Errorf("failed to execute scene template: %w", err)
	}

	return nil
}
