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
	"encoding/xml"
	"errors"
	"path/filepath"
)

const (
	tileHorizontalFlipMask = 0x80000000
	tileVerticalFlipMask   = 0x40000000
	tileDiagonalFlipMask   = 0x20000000
	tileFlip               = tileHorizontalFlipMask | tileVerticalFlipMask | tileDiagonalFlipMask
	tileGIDMask            = 0x0fffffff
)

var (
	// ErrInvalidTileGID error is returned when tile GID is not found
	ErrInvalidTileGID = errors.New("tiled: invalid tile GID")
)

// Map contains three different kinds of layers.
// Tile layers were once the only type, and are simply called layer, object layers have the objectgroup tag
// and image layers use the imagelayer tag. The order in which these layers appear is the order in which the
// layers are rendered by Tiled
type Map struct {
	// Loader for loading additional data
	loader *Loader
	// Base directory for loading additional data
	baseDir string

	// The TMX format version, generally 1.0.
	Version string `xml:"title,attr"`
	// Map orientation. Tiled supports "orthogonal", "isometric", "staggered" (since 0.9) and "hexagonal" (since 0.11).
	Orientation string `xml:"orientation,attr"`
	// The order in which tiles on tile layers are rendered. Valid values are right-down (the default), right-up, left-down and left-up.
	// In all cases, the map is drawn row-by-row. (since 0.10, but only supported for orthogonal maps at the moment)
	RenderOrder string `xml:"renderorder,attr"`
	// The map width in tiles.
	Width int `xml:"width,attr"`
	// The map height in tiles.
	Height int `xml:"height,attr"`
	// The width of a tile.
	TileWidth int `xml:"tilewidth,attr"`
	// The height of a tile.
	TileHeight int `xml:"tileheight,attr"`
	// Only for hexagonal maps. Determines the width or height (depending on the staggered axis) of the tile's edge, in pixels.
	HexSideLength int `xml:"hexsidelength,attr"`
	// For staggered and hexagonal maps, determines which axis ("x" or "y") is staggered. (since 0.11)
	StaggerAxis int `xml:"staggeraxis,attr"`
	// For staggered and hexagonal maps, determines whether the "even" or "odd" indexes along the staggered axis are shifted. (since 0.11)
	StaggerIndex int `xml:"staggerindex,attr"`
	// The background color of the map. (since 0.9, optional, may include alpha value since 0.15 in the form #AARRGGBB)
	BackgroundColor string `xml:"backgroundcolor,attr"`
	// Stores the next available ID for new objects. This number is stored to prevent reuse of the same ID after objects have been removed. (since 0.11)
	NextObjectID uint32 `xml:"nextobjectid,attr"`
	// Custom properties
	Properties *Properties `xml:"properties>property"`
	// Map tilesets
	Tilesets []*Tileset `xml:"tileset"`
	// Map layers
	Layers []*Layer `xml:"layer"`
	// Map object groups
	ObjectGroups []*ObjectGroup `xml:"objectgroup"`
	// Image layers
	ImageLayers []*ImageLayer `xml:"imagelayer"`
	// Group layers
	Groups []*Group `xml:"group"`
}

func (m *Map) initTileset(ts *Tileset) (*Tileset, error) {
	if ts.SourceLoaded {
		return ts, nil
	}
	if len(ts.Source) == 0 {
		ts.baseDir = m.baseDir
		ts.SourceLoaded = true
		return ts, nil
	}
	sourcePath := m.GetFileFullPath(ts.Source)
	f, err := m.loader.open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := xml.NewDecoder(f)

	tse := &Tileset{}
	if err := d.Decode(tse); err != nil {
		return nil, err
	}

	tse.baseDir = filepath.Dir(sourcePath)
	tse.Source = ts.Source
	tse.SourceLoaded = true
	tse.FirstGID = ts.FirstGID

	return tse, nil
}

// TileGIDToTile is used to find tile data by GID
func (m *Map) TileGIDToTile(gid uint32) (*LayerTile, error) {
	if gid == 0 {
		return NilLayerTile, nil
	}

	gidBare := gid &^ tileFlip

	for i := len(m.Tilesets) - 1; i >= 0; i-- {
		if m.Tilesets[i].FirstGID <= gidBare {
			ts, err := m.initTileset(m.Tilesets[i])
			if err != nil {
				return nil, err
			}
			return &LayerTile{
				ID:             gidBare - ts.FirstGID,
				Tileset:        ts,
				HorizontalFlip: gid&tileHorizontalFlipMask != 0,
				VerticalFlip:   gid&tileVerticalFlipMask != 0,
				DiagonalFlip:   gid&tileDiagonalFlipMask != 0,
				Nil:            false,
			}, nil
		}
	}

	return nil, ErrInvalidTileGID
}

// GetFileFullPath returns path to file relative to map file
func (m *Map) GetFileFullPath(fileName string) string {
	return filepath.Join(m.baseDir, fileName)
}

// UnmarshalXML decodes a single XML element beginning with the given start element.
func (m *Map) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias Map

	item := Alias{
		loader:      m.loader,
		baseDir:     m.baseDir,
		RenderOrder: "right-down",
	}

	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}

	// Decode layers data
	for i := 0; i < len(item.Layers); i++ {
		l := item.Layers[i]
		if err := l.DecodeLayer((*Map)(&item)); err != nil {
			return err
		}
	}

	*m = (Map)(item)
	return nil
}
