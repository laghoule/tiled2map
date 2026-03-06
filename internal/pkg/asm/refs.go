package asm

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const (
	tileAttribute = "attr"
	tilHeaderSize = 3
)

// TileRefsTemplateData is the data structure for the tiles references template
type TileRefsTemplateData struct {
	GID         int
	SourceImage string
	Attribute   string
}

// TileRefsTemplatePayload is the top-level data passed to the tiles references template
type TileRefsTemplatePayload struct {
	Prefix     string
	BufferSize int
	TilesRefs  []TileRefsTemplateData
}

// createTilesRefs generates the ASM tile references file
func (a *ASMLinker) createTilesRefs() error {
	filename := filepath.Join(a.FileOutput.Path, fmt.Sprintf("%s-refs.inc", a.FileOutput.FilePrefix))
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
					switch val := prop.Value.(type) {
					case float64:
						attribute = fmt.Sprintf("%08bb", int(val))
					case string:
						validASMLabel := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
						upper := strings.ToUpper(val)
						if !validASMLabel.MatchString(upper) {
							return fmt.Errorf("invalid attribute label for tile %d: %q (must match ^[A-Za-z_][A-Za-z0-9_]*$)", tileInfo.GID, val)
						}
						attribute = upper
					default:
						return fmt.Errorf("invalid attribute value for tile %d: %v", tileInfo.GID, prop.Value)
					}
					break
				}
			}
		}

		// If no attribute is set, default to 00000000b
		if attribute == "" {
			attribute = "00000000b"
		}

		tilesRefs = append(tilesRefs, TileRefsTemplateData{
			GID:         tileInfo.GID,
			SourceImage: tileInfo.SourceImage,
			Attribute:   attribute,
		})
	}

	bufferSize := (a.TileMap.TileWidth * a.TileMap.TileHeight * len(a.TilesInfo)) + tilHeaderSize

	payload := TileRefsTemplatePayload{
		Prefix:     a.FileOutput.FilePrefix,
		BufferSize: bufferSize,
		TilesRefs:  tilesRefs,
	}

	tpl, err := template.ParseFiles("tmpl/tiles_props.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	err = tpl.Execute(asmFile, payload)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}
