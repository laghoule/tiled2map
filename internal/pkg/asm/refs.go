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
)

// TileRefsData is the data structure for the tiles references template
type TileRefsData struct {
	GID         int
	SourceImage string
	Attribute   string
}

// TileData is the data structure for the tiles data template
type TileData struct {
	Width  int
	Height int
	Count  int
}

// TileRefsTemplatePayload is the top-level data passed to the tiles references template
type TileRefsTemplatePayload struct {
	Prefix     string
	BufferSize int

	TileData  TileData
	TilesRefs []TileRefsData
}

// createTilesRefs generates the ASM tile references file
func (a *ASMLinker) createTilesRefs() error {
	filename := filepath.Join(a.FileOutput.Path, fmt.Sprintf("%s-refs.inc", a.FileOutput.FilePrefix))
	asmFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer asmFile.Close()

	tilesRefs := []TileRefsData{}

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

		tilesRefs = append(tilesRefs, TileRefsData{
			GID:         tileInfo.GID,
			SourceImage: tileInfo.SourceImage,
			Attribute:   attribute,
		})
	}

	tilesCount := len(a.TilesInfo)
	bufferSize := (a.TileMap.TileWidth * a.TileMap.TileHeight * tilesCount)

	payload := TileRefsTemplatePayload{
		Prefix:     a.FileOutput.FilePrefix,
		BufferSize: bufferSize,
		TileData: TileData{
			Width:  a.TileMap.TileWidth,
			Height: a.TileMap.TileHeight,
			Count:  tilesCount,
		},
		TilesRefs: tilesRefs,
	}

	funcMap := template.FuncMap{
		"toUpper": strings.ToUpper,
	}

	tpl, err := template.New("tiles_props.tmpl").Funcs(funcMap).ParseFS(tmplFS, "tmpl/tiles_props.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	err = tpl.Execute(asmFile, payload)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}
