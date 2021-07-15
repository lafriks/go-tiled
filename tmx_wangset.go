package tiled

import (
	"errors"
	"strconv"
	"strings"
)

// https://doc.mapeditor.org/en/stable/reference/tmx-map-format/#wangsets

// Contains the list of Wang sets defined for this tileset.
// Can contain any number: <wangset>
type WangSets []*WangSet

// Defines a list of corner colors and a list of edge colors, and any number of Wang tiles using these colors.
type WangSet struct {
	Name       string       `xml:"name,attr"` // The name of the Wang set.
	Type       string       `xml:"type,attr"` // ex. corner
	TileId     uint32       `xml:"tile,attr"` // The tile ID of the tile representing this Wang set.
	WangColors []*WangColor `xml:"wangcolor"`
	WangTiles  []*WangTile  `xml:"wangtile"`
}

//A color that can be used to define the corner and/or edge of a Wang tile.
type WangColor struct {
	Name        string  `xml:"name,attr"`        //  The name of this color.
	Color       string  `xml:"color,attr"`       // The color in #RRGGBB format (example: #c17d11).
	TileId      uint32  `xml:"tile,attr"`        // The tile ID of the tile representing this color.
	Probability float32 `xml:"probability,attr"` // The relative probability that this color is chosen over others in case of multiple options. (defaults to 0)
}

//Defines a Wang tile, by referring to a tile in the tileset and associating it with a certain Wang ID.
type WangTile struct {
	TileId uint32 `xml:"tileid,attr"` // The tile ID.

	// The Wang ID, given by a comma-separated list of indexes (starting from 1, because 0 means _unset_)
	// referring to the Wang colors in the Wang set in the following order:
	// top, top right, right, bottom right, bottom, bottom left, left, top left (since Tiled 1.5).
	// Before Tiled 1.5, the Wang ID was saved as a 32-bit unsigned integer stored in the format
	// 0xCECECECE (where each C is a corner color and each E is an edge color, in reverse order).
	WangId string `xml:"wangid,attr"` //
}

// getWangColors returns the wangcolors for the tileId. If corner type is used it will return an array of len 4
// topRight, bottomRight, bottom left, top left
// if corner type is not used it will return an array of len 8 in the following order.
// top, top right, right, bottom right, bottom, bottom left, left, top left
// if there is no wangcolor assigned to a part of the tile it will return an nil pointer instead for that index
func (w *WangSet) GetWangColors(tileId uint32) ([]*WangColor, error) {

	if w.WangColors == nil {
		return nil, errors.New("no wangcolors found on this wangset")
	}

	var tile *WangTile
	for _, t := range w.WangTiles {
		if t.TileId == tileId {
			tile = t
			break
		}
	}
	if tile == nil {
		return nil, errors.New("no wangtile matches the given Id")
	}

	// convert from CSV to array of strings
	wangIdsString := strings.Split(tile.WangId, ",")

	// convert from array of strings to slice of uint32
	var wangIds []uint32 // will contain a slice of the wangIds
	for _, v := range wangIdsString {
		id64, err := strconv.ParseUint(v, 10, 32)

		if err != nil {
			return nil, errors.New("internal error")
		}

		// uint64 to uint32
		id := uint32(id64)

		wangIds = append(wangIds, id)
	}

	var wangColors []*WangColor

	if w.Type == "corner" { // missing top, right, bottom and left..

		if len(w.WangColors) < len(wangIds)/2 {
			return nil, errors.New("too few wangcolors found")
		}

		for i, id := range wangIds {
			if i%2 == 0 { // skip even indices
				continue
			}

			if id == 0 { // no color assigned if id is 0, set to nil
				wangColors = append(wangColors, nil)
			} else {

				wangColors = append(wangColors, w.WangColors[id-1]) // minus 1 because there is no 0 id, since 0 means unassigned
			}
		}

	} else { // type != corner

		if len(w.WangColors) < len(wangIds) {
			return nil, errors.New("too few wangcolors found")
		}

		for _, id := range wangIds {

			if id == 0 { // no color assigned if id is 0, set to nil
				wangColors = append(wangColors, nil)
			} else {

				wangColors = append(wangColors, w.WangColors[id-1]) // minus 1 because there is no 0 id, since 0 means unassigned
			}
		}
	}

	return wangColors, nil

}
