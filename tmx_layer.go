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
)

// NilLayerTile is reusable layer tile that is nil
var NilLayerTile = &LayerTile{Nil: true}

var (
	// ErrInvalidDecodedTileCount error is returned when layer data has invalid tile count
	ErrInvalidDecodedTileCount = errors.New("tiled: invalid decoded tile count")
	// ErrEmptyLayerData error is returned when layer contains no data
	ErrEmptyLayerData = errors.New("tiled: missing layer data")
	// ErrUnknownEncoding error is returned when kayer data has unknown encoding
	ErrUnknownEncoding = errors.New("tiled: unknown data encoding")
)

// LayerTile is a layer tile
type LayerTile struct {
	// Tile ID
	ID uint32
	// Tile tileset
	Tileset *Tileset
	// Horizontal tile image flip
	HorizontalFlip bool
	// Vertical tile image flip
	VerticalFlip bool
	// Diagonal tile image flip
	DiagonalFlip bool
	// Tile is nil
	Nil bool
}

// IsNil returs if tile is nil
func (t *LayerTile) IsNil() bool {
	return t.Nil
}

// Layer is a map layer
type Layer struct {
	_map *Map
	// Unique ID of the layer.
	// Each layer that added to a map gets a unique id. Even if a layer is deleted,
	// no layer ever gets the same ID. Can not be changed in Tiled. (since Tiled 1.2)
	ID uint32 `xml:"id,attr"`
	// The name of the layer.
	Name string `xml:"name,attr"`
	// The opacity of the layer as a value from 0 to 1. Defaults to 1.
	Opacity float32 `xml:"opacity,attr"`
	// Whether the layer is shown (1) or hidden (0). Defaults to 1.
	Visible bool `xml:"visible,attr"`
	// Rendering offset for this layer in pixels. Defaults to 0. (since 0.14)
	OffsetX int `xml:"offsetx,attr"`
	// Rendering offset for this layer in pixels. Defaults to 0. (since 0.14)
	OffsetY int `xml:"offsety,attr"`
	// Custom properties
	Properties Properties `xml:"properties>property"`
	// This is the attribute you'd like to use, not Data. Tile entry at (x,y) is obtained using l.DecodedTiles[y*map.Width+x].
	Tiles []*LayerTile
	// Data
	data *Data
	// Tileset if only one is used in layer
	tileset *Tileset
	// Set when all entries of the layer are NilTile
	empty bool
}

// IsEmpty checks if layer has tiles other than nil
func (l *Layer) IsEmpty() bool {
	return l.empty
}

func (l *Layer) decodeLayerXML() (gids []uint32, err error) {
	if len(l.data.DataTiles) != l._map.Width*l._map.Height {
		return []uint32{}, ErrInvalidDecodedTileCount
	}

	gids = make([]uint32, len(l.data.DataTiles))
	for i := 0; i < len(gids); i++ {
		gids[i] = l.data.DataTiles[i].GID
	}

	return gids, nil
}

func (l *Layer) decodeLayerCSV() ([]uint32, error) {
	gids, err := l.data.decodeCSV()
	if err != nil {
		return []uint32{}, err
	}

	if len(gids) != l._map.Width*l._map.Height {
		return []uint32{}, ErrInvalidDecodedTileCount
	}

	return gids, nil
}

func (l *Layer) decodeLayerBase64() ([]uint32, error) {
	dataBytes, err := l.data.decodeBase64()
	if err != nil {
		return []uint32{}, err
	}

	if len(dataBytes) != l._map.Width*l._map.Height*4 {
		return []uint32{}, ErrInvalidDecodedTileCount
	}

	gids := make([]uint32, l._map.Width*l._map.Height)

	j := 0
	for y := 0; y < l._map.Height; y++ {
		for x := 0; x < l._map.Width; x++ {
			gid := uint32(dataBytes[j]) +
				uint32(dataBytes[j+1])<<8 +
				uint32(dataBytes[j+2])<<16 +
				uint32(dataBytes[j+3])<<24
			j += 4

			gids[y*l._map.Width+x] = gid
		}
	}

	return gids, nil
}

func (l *Layer) decodeTiles() error {
	var gids []uint32
	var err error
	switch l.data.Encoding {
	case "csv":
		if gids, err = l.decodeLayerCSV(); err != nil {
			return err
		}
	case "base64":
		if gids, err = l.decodeLayerBase64(); err != nil {
			return err
		}
	case "": // XML "encoding"
		if gids, err = l.decodeLayerXML(); err != nil {
			return err
		}
	default:
		return ErrUnknownEncoding
	}

	l.Tiles = make([]*LayerTile, len(gids))
	for j := 0; j < len(l.Tiles); j++ {
		l.Tiles[j], err = l._map.TileGIDToTile(gids[j])
		if err != nil {
			return err
		}
	}

	return nil
}

// DecodeLayer decodes layer data
func (l *Layer) DecodeLayer(m *Map) error {
	l._map = m
	if l.data == nil {
		return ErrEmptyLayerData
	}

	if err := l.decodeTiles(); err != nil {
		return err
	}

	// Data is not needed anymore
	l.data = nil

	var tileset *Tileset
	for i := 0; i < len(l.Tiles); i++ {
		tile := l.Tiles[i]
		if !tile.Nil {
			if tileset == nil {
				tileset = tile.Tileset
			} else if tileset != tile.Tileset {
				l.tileset = nil
				l.empty = false
				return nil
			}
		}
	}

	l.tileset = tileset
	l.empty = false

	if tileset == nil {
		l.empty = true
	}

	return nil
}

// UnmarshalXML decodes a single XML element beginning with the given start element.
func (l *Layer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type InternalAlias Layer

	type Alias struct {
		InternalAlias
		// Layer data in raw format
		Data *Data `xml:"data"`
	}

	item := Alias{InternalAlias: InternalAlias{
		Opacity: 1,
		Visible: true,
	}}

	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}

	*l = (Layer)(item.InternalAlias)

	l.data = item.Data

	return nil
}
