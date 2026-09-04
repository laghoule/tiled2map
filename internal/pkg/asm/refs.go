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
	tileAttribute         = "attr"
	tileAttributeTmplFile = "tmpl/tiles_attributes.tmpl"
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
	Prefix             string
	BufferSize         int
	GenerateAttributes bool

	TileData  TileData
	TilesRefs []TileRefsData
}

// createTilesRefs generates the ASM tile references file
func (a *ASMLinker) createTilesRefs() error {
	filename := filepath.Join(a.fileOutput.Path, fmt.Sprintf("%s-refs.inc", a.fileOutput.FilePrefix))
	asmFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer asmFile.Close()

	tilesRefs := []TileRefsData{}

	if a.generateAttributes {
		for _, tileInfo := range a.tilesInfo {
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
	}

	tilesCount := len(a.tilesInfo)
	bufferSize := (a.tileMap.TileWidth * a.tileMap.TileHeight * tilesCount)

	payload := TileRefsTemplatePayload{
		Prefix:        a.fileOutput.FilePrefix,
		BufferSize:    bufferSize,
		GenerateAttributes: a.generateAttributes,
		TileData: TileData{
			Width:  a.tileMap.TileWidth,
			Height: a.tileMap.TileHeight,
			Count:  tilesCount,
		},
		TilesRefs: tilesRefs,
	}

	funcMap := template.FuncMap{
		"toUpper": strings.ToUpper,
	}

	tpl, err := template.New("tiles_attributes.tmpl").Funcs(funcMap).ParseFS(tmplFS, tileAttributeTmplFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	err = tpl.Execute(asmFile, payload)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}
