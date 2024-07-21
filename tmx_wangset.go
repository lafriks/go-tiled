package tiled

import (
	"errors"
	"strconv"
	"strings"
)

// WangSets contains the list of Wang sets defined for this tileset.
// https://doc.mapeditor.org/en/stable/reference/tmx-map-format/#wangsets
// Can contain any number: <wangset>
type WangSets []*WangSet

// WangSet defines a list of corner colors and a list of edge colors, and any number of Wang tiles using these colors.
type WangSet struct {
	// The name of the Wang set.
	Name string `xml:"name,attr"`
	// The class of the Wang set. (renamed from 'type' since 1.9)
	Class string `xml:"class,attr"`
	// WangSet type.
	//
	// Deprecated: replaced by Class since 1.9
	Type string `xml:"type,attr"`
	// The tile ID of the tile representing this Wang set.
	TileID int64 `xml:"tile,attr"`
	// The list of corner and/or edge colors.
	WangColors []*WangColor `xml:"wangcolor"`
	// The list of wang tiles.
	WangTiles []*WangTile `xml:"wangtile"`
}

// WangColor that can be used to define the corner and/or edge of a Wang tile.
type WangColor struct {
	// The name of this color.
	Name string `xml:"name,attr"`
	// The class of this color. (since 1.9)
	Class string `xml:"class,attr"`
	// The color in #RRGGBB format (example: #c17d11).
	Color string `xml:"color,attr"`
	// The tile ID of the tile representing this color.
	TileID int64 `xml:"tile,attr"`
	// The relative probability that this color is chosen over others in case of multiple options. (defaults to 0)
	Probability float32 `xml:"probability,attr"`
}

// WangTile by referring to a tile in the tileset and associating it with a certain Wang ID.
type WangTile struct {
	// The tile ID.
	TileID uint32 `xml:"tileid,attr"`
	// WangID, given by a comma-separated list of indexes (starting from 1, because 0 means _unset_)
	// referring to the Wang colors in the Wang set in the following order:
	// top, top right, right, bottom right, bottom, bottom left, left, top left (since Tiled 1.5).
	// Before Tiled 1.5, the Wang ID was saved as a 32-bit unsigned integer stored in the format
	// 0xCECECECE (where each C is a corner color and each E is an edge color, in reverse order).
	WangID string `xml:"wangid,attr"`
}

// WangPosition Wang Color mapping to position
type WangPosition int

const (
	// Top represents the top part of the tile
	Top WangPosition = iota
	// TopRight represents the top right part of the tile
	TopRight
	// Right represents the right part of the tile
	Right
	// BottomRight represents the bottom right part of the tile
	BottomRight
	// Bottom represents the bottom part of the tile
	Bottom
	// BottomLeft represents the bottom left part of the tile
	BottomLeft
	// Left represents the left part of the tile
	Left
	// TopLeft represents the top left part of the tile
	TopLeft
)

// GetWangColors returns the wangcolors for the tileId. If corner type is used it will return an array of len 4
// topRight, bottomRight, bottom left, top left
// if corner type is not used it will return an array of len 8 in the following order.
// top, top right, right, bottom right, bottom, bottom left, left, top left
// if there is no wangcolor assigned to a part of the tile it will return an nil pointer instead for that index
func (w *WangSet) GetWangColors(tileID uint32) (map[WangPosition]*WangColor, error) {
	if w.WangColors == nil {
		return nil, errors.New("no wangcolors found on this wangset")
	}

	var tile *WangTile
	for _, t := range w.WangTiles {
		if t.TileID == tileID {
			tile = t
			break
		}
	}
	if tile == nil {
		return nil, errors.New("no wangtile matches the given Id")
	}

	// convert from CSV to array of strings
	wangIDsString := strings.Split(tile.WangID, ",")

	// convert from array of strings to slice of uint32
	var wangIDs []uint32 // will contain a slice of the wangIDs
	for _, v := range wangIDsString {
		id64, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return nil, errors.New("internal error")
		}

		// uint64 to uint32
		id := uint32(id64)

		wangIDs = append(wangIDs, id)
	}

	wangColors := make(map[WangPosition]*WangColor)

	for i, id := range wangIDs {
		if id == 0 { // no color assigned if id is 0, set to nil
			wangColors[WangPosition(i)] = nil
		} else {
			wangColors[WangPosition(i)] = w.WangColors[id-1]
		}
	}

	return wangColors, nil
}
