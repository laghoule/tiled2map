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
func createScene(filePrefix string, d dimension) error {
	tpl, err := template.ParseFiles("tmpl/scene.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse scene template: %w", err)
	}

	scene := []SceneTemplateData{}

	var north, south, east, west string
	for y := range d.height {
		for x := range d.width {

			// FIXME: not working
			
			if y == 0 {
				north = "0"
			} else {
				north = fmt.Sprintf("SCENE_%d_%d", x, y-1)
			}
			if y == d.height-1 {
				south = "0"
			} else {
				south = fmt.Sprintf("SCENE_%d_%d", x, y+1)
			}

			if x == d.width-1 {
				east = "0"
			} else {
				east = fmt.Sprintf("SCENE_%d_%d", x+1, y)
			}
			if x == 0 {
				west = "0"
			} else {
				west = fmt.Sprintf("SCENE_%d_%d", x-1, y)
			}

			scene = append(scene, SceneTemplateData{
				Name:      fmt.Sprintf("Scene_%d_%d", x, y),
				BGOffset:  int(x) * 16, // FIXME (16)
				FGOffset:  int(y) * 16, // FIXME (16)
				NorthName: north,
				SouthName: south,
				EastName:  east,
				WestName:  west,
				MusicName: fmt.Sprintf("Music_%d_%d", x, y),
			})
		}
	}
	
	filename := fmt.Sprintf("%s-scene.%s", filePrefix, includeExt)
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
