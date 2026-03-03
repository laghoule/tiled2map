package asm

import (
	"fmt"
	"os"
	"text/template"
)

const (
	tileAttribute = "attr"
)

// TileRefsTemplateData is the data structure for the tiles references template
type TileRefsTemplateData struct {
	GID         int
	SourceImage string
	Attribute   string
}

// createTilesRefs generates the ASM tile references file
func (a *ASMLinker) createTilesRefs() error {
	filename := fmt.Sprintf("%s-refs.inc", a.FilePrefix)
	asmFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer asmFile.Close()

	tilesRefs := []TileRefsTemplateData{}

	for _, tileInfo := range a.TilesInfo {
		attribute := ""
		for _, tile := range tileInfo.Tiles {
			for _, prop := range tile.Properties {
				if prop.Name == tileAttribute {
					val, ok := prop.Value.(float64)
					if !ok {
						return fmt.Errorf("invalid attribute value for tile %d: %v", tileInfo.GID, prop.Value)
					}
					attribute = fmt.Sprintf("%08bb", int(val))
					break
				}
			}
		}

		if attribute == "" {
			attribute = "00000000b"
		}

		tilesRefs = append(tilesRefs, TileRefsTemplateData{
			GID:         tileInfo.GID,
			SourceImage: tileInfo.SourceImage,
			Attribute:   attribute,
		})
	}

	tpl, err := template.ParseFiles("tmpl/tiles_props.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	err = tpl.Execute(asmFile, tilesRefs)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}
