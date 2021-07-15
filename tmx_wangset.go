package tiled

// Properties wraps any number of custom properties
type WangSets []*WangSet

type WangSet struct {
	Name       string       `xml:"name,attr"`
	Type       string       `xml:"type,attr"`
	WangColors []*WangColor `xml:"wangcolor"`
	WangTiles  []*WangTile  `xml:"wangtile"`
}

type WangColor struct {
	Name  string `xml:"name,attr"`
	Color string `xml:"color,attr"`
}

type WangTile struct {
	TileId string `xml:"tileid,attr"`
	WangId string `xml:"wangid,attr"`
}
