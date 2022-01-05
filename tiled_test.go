/*
Copyright (c) 2017 Lauris Bukšis-Haberkorns <lauris@nix.lv>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package tiled

import (
	"bytes"
	"embed"
	"encoding/xml"
	"image/color"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetAssetsDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "assets")
}

func TestLoadFromReader(t *testing.T) {
	r := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<map version="1.2" tiledversion="1.2.1" orientation="orthogonal" renderorder="right-down" width="4" height="4" tilewidth="16" tileheight="16" infinite="0" nextlayerid="2" nextobjectid="2">
<layer id="1" name="Tile Layer 1" width="4" height="4">
<data encoding="csv">
0,0,0,0,
0,0,0,0,
0,0,0,0,
0,0,0,0
</data>
</layer>
</map>`)
	m, err := LoadFromReader(GetAssetsDirectory(), r)

	assert.NoError(t, err)
	assert.NotNil(t, m)
	assert.Len(t, m.Layers, 1)
}

func TestLoadFromReaderError(t *testing.T) {
	r := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<map version="1.2" tiledversion="1.2.1" orientation="orthogonal" renderorder="right-down" width="4" height="4" tilewidth="16" tileheight="16" infinite="0" nextlayerid="2" nextobjectid="2">
<layer id="1" name="Tile Layer 1" width="4" height="4">
<data encoding="csv">
</layer>
</map>`)
	m, err := LoadFromReader(GetAssetsDirectory(), r)

	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestLoadFromFile(t *testing.T) {
	m, err := LoadFromFile(filepath.Join(GetAssetsDirectory(), "test.tmx"))

	assert.NoError(t, err)
	assert.NotNil(t, m)

	assert.Len(t, m.Layers, 1)
	assert.Equal(t, uint32(1), m.Layers[0].ID)

	// Test ObjectGroups.Visible defaults to true
	assert.Len(t, m.ObjectGroups, 1)
	assert.Equal(t, uint32(2), m.ObjectGroups[0].ID)
	assert.Equal(t, true, *m.ObjectGroups[0].Visible)

	// Test Object.Visible defaults to true
	assert.Len(t, m.ObjectGroups[0].Objects, 1)
	assert.Equal(t, uint32(2), m.ObjectGroups[0].Objects[0].ID)
	assert.Equal(t, true, *m.ObjectGroups[0].Objects[0].Visible)
}

func TestLoadFromFileError(t *testing.T) {
	m, err := LoadFromFile(filepath.Join(GetAssetsDirectory(), "invalid.tmx"))

	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestExternalTilesetImageLoaded(t *testing.T) {
	m, err := LoadFromFile(filepath.Join(GetAssetsDirectory(), "test2.tmx"))

	assert.NoError(t, err)
	assert.NotNil(t, m)

	for _, layer := range m.Layers {
		var tileset *Tileset
		for _, tile := range layer.Tiles {
			if tile.ID > 0 {
				tileset = tile.Tileset
				assert.NotNil(t, tileset)
				if tileset != nil {
					assert.NotNil(t, tileset.Image)
					assert.Equal(t, "ProjectUtumno_full.png", tileset.Image.Source)
					assert.Contains(t, m.Tilesets, tileset)
				}
			}
		}
	}
}

func TestImageLayer(t *testing.T) {
	m, err := LoadFromFile(filepath.Join(GetAssetsDirectory(), "imagelayer.tmx"))

	assert.NoError(t, err)
	assert.NotNil(t, m)

	assert.Len(t, m.ImageLayers, 1)

	layer := m.ImageLayers[0]
	assert.Equal(t, uint32(2), layer.ID)
	assert.Equal(t, "Image Layer", layer.Name)
	assert.Equal(t, 0, layer.OffsetX)
	assert.Equal(t, 0, layer.OffsetY)
	assert.Equal(t, float32(1.0), layer.Opacity)
	assert.Equal(t, true, layer.Visible)

	image := layer.Image
	assert.NotNil(t, image)

	assert.Equal(t, image.Source, "background.jpg", image.Source)
}

func TestGroup(t *testing.T) {
	m, err := LoadFromFile(filepath.Join(GetAssetsDirectory(), "groups.tmx"))

	assert.NoError(t, err)
	assert.NotNil(t, m)

	assert.Len(t, m.Layers, 1)
	assert.Len(t, m.Groups, 1)

	a := m.Groups[0]
	assert.Equal(t, uint32(2), a.ID)
	assert.Equal(t, "A", a.Name)
	assert.Len(t, a.ImageLayers, 1)
	assert.Len(t, a.Groups, 1)

	b := a.Groups[0]
	assert.Equal(t, uint32(4), b.ID)
	assert.Equal(t, "B", b.Name)
	assert.Len(t, b.Layers, 1)
	assert.Len(t, b.Groups, 1)

	bL := b.Layers[0]
	assert.Equal(t, uint32(3), bL.ID)
	assert.Equal(t, bL.Name, "Tile Layer 2")
	assert.Len(t, bL.Tiles, 400)

	c := b.Groups[0]
	assert.Equal(t, uint32(8), c.ID)
	assert.Equal(t, "C", c.Name)
	assert.Len(t, c.ObjectGroups, 1)
	assert.Len(t, c.Groups, 0)
}

func TestFont(t *testing.T) {
	m, err := LoadFromFile(filepath.Join(GetAssetsDirectory(), "font.tmx"))

	assert.NoError(t, err)
	assert.NotNil(t, m)

	if assert.Len(t, m.ObjectGroups, 1) {
		if assert.Len(t, m.ObjectGroups[0].Objects, 1) {
			if assert.NotNil(t, m.ObjectGroups[0].Objects[0].Text) {
				text := m.ObjectGroups[0].Objects[0].Text

				assert.Equal(t, "Hello World", text.Text)
				assert.Equal(t, "sans-serif", text.FontFamily)
				assert.Equal(t, 16, text.Size)
				assert.Equal(t, true, text.Wrap)
				assert.Equal(t, &HexColor{}, text.Color)
				assert.Equal(t, false, text.Bold)
				assert.Equal(t, false, text.Italic)
				assert.Equal(t, false, text.Underline)
				assert.Equal(t, false, text.Strikethrough)
				assert.Equal(t, true, *text.Kerning)
				assert.Equal(t, "left", text.HAlign)
				assert.Equal(t, "top", text.VAlign)
			}
		}
	}
}

type testFileSystem struct {
	AttemptedOpen []string
}

func (t *testFileSystem) Open(filename string) (fs.File, error) {
	t.AttemptedOpen = append(t.AttemptedOpen, filename)
	if filepath.Base(filename) == "loader.tmx" {
		return os.Open(filepath.Join(GetAssetsDirectory(), "loader.tmx"))
	}
	return nil, os.ErrNotExist
}

func TestLoader(t *testing.T) {
	fs := &testFileSystem{}
	loader := &Loader{
		FileSystem: fs,
	}

	mapFile := filepath.Join(GetAssetsDirectory(), "loader.tmx")
	m, err := loader.LoadFromFile(mapFile)

	if assert.Error(t, err) {
		assert.True(t, os.IsNotExist(err), "expecting no exist error")
	}
	assert.Nil(t, m)

	assert.Equal(t, []string{mapFile, filepath.Join(GetAssetsDirectory(), "..", "README.md")}, fs.AttemptedOpen)
}

//go:embed assets/**
var assetsFS embed.FS

func TestEmbeddedLoader(t *testing.T) {
	loader := &Loader{
		FileSystem: assetsFS,
	}
	tcs := []struct {
		name string
		load func() (*Map, error)
	}{
		{
			name: "LoadFromReader",
			load: func() (*Map, error) {
				file, err := assetsFS.Open("assets/test2.tmx")
				if err != nil {
					return nil, err
				}
				return loader.LoadFromReader("assets", file)
			},
		},
		{
			name: "LoadFromFile",
			load: func() (*Map, error) {
				return loader.LoadFromFile("assets/test2.tmx")
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			m, err := tc.load()
			assert.NoError(t, err)

			assert.Equal(t, m.Version, "1.2")
			assert.Equal(t, m.TiledVersion, "1.2.3")

			assert.Len(t, m.Tilesets, 1)
			tileset := m.Tilesets[0]
			assert.Equal(t, tileset.Version, "1.2")
			assert.Equal(t, tileset.TiledVersion, "1.2.3")
		})
	}
}

func TestParseHexColor(t *testing.T) {
	type test struct {
		name  string
		attr  xml.Attr
		color color.RGBA
	}
	cases := []test{
		{
			name: "Wihout alpha",
			attr: xml.Attr{
				Value: "#aabbcc",
			},
			color: color.RGBA{
				R: 170,
				G: 187,
				B: 204,
				A: 255,
			},
		},
		{
			name: "With alpha",
			attr: xml.Attr{
				Value: "#aabbccdd",
			},
			color: color.RGBA{
				R: 187,
				G: 204,
				B: 221,
				A: 170,
			},
		},
		{
			name: "Simplified",
			attr: xml.Attr{
				Value: "#fff",
			},
			color: color.RGBA{
				R: 255,
				G: 255,
				B: 255,
				A: 255,
			},
		},
		{
			name: "Simplified with alpha",
			attr: xml.Attr{
				Value: "#ffff",
			},
			color: color.RGBA{
				R: 255,
				G: 255,
				B: 255,
				A: 255,
			},
		},
		{
			name: "Different values",
			attr: xml.Attr{
				Value: "#010203",
			},
			color: color.RGBA{
				R: 1,
				G: 2,
				B: 3,
				A: 255,
			},
		},
		{
			name: "uppercase",
			attr: xml.Attr{
				Value: "#ABCDEF",
			},
			color: color.RGBA{
				R: 171,
				G: 205,
				B: 239,
				A: 255,
			},
		},
	}

	for _, c := range cases {
		color := &HexColor{}
		color.UnmarshalXMLAttr(c.attr)
		assert.Equal(t, c.color, color.c, c.name)
	}
}

func TestFormatHexColor(t *testing.T) {
	type test struct {
		name  string
		color HexColor
		hex   string
	}
	cases := []test{
		{
			name:  "alpha ignored at 255",
			color: NewHexColor(255, 255, 255, 255),
			hex:   "#ffffff",
		},
		{
			name:  "transparent black",
			color: NewHexColor(255, 255, 255, 0),
			hex:   "#00ffffff",
		},
		{
			name:  "50% transparency",
			color: NewHexColor(255, 255, 255, 128),
			hex:   "#80ffffff",
		},
		{
			name:  "values encoded with 2 numbers",
			color: NewHexColor(16, 32, 48, 255),
			hex:   "#102030",
		},
		{
			name:  "values encoded with 1 numbers",
			color: NewHexColor(1, 2, 3, 255),
			hex:   "#010203",
		},
		{
			name:  "zero",
			color: NewHexColor(0, 0, 0, 0),
			hex:   "#00000000",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.hex, c.color.String(), c.name)
	}
}

func TestVersions(t *testing.T) {
	m, err := LoadFromFile(filepath.Join(GetAssetsDirectory(), "test2.tmx"))

	assert.Nil(t, err)

	assert.Equal(t, m.Version, "1.2")
	assert.Equal(t, m.TiledVersion, "1.2.3")

	assert.Len(t, m.Tilesets, 1)

	tileset := m.Tilesets[0]

	assert.Equal(t, tileset.Version, "1.2")
	assert.Equal(t, tileset.TiledVersion, "1.2.3")
}
