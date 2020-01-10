package tiled

import (
	"path/filepath"
)

// Tileset is collection of tiles
type Tileset struct {
	// Base directory
	baseDir string

	// The first global tile ID of this tileset (this global ID maps to the first tile in this tileset).
	FirstGID uint32 `xml:"firstgid,attr"`
	// If this tileset is stored in an external TSX (Tile Set XML) file, this attribute refers to that file.
	// That TSX file has the same structure as the <tileset> element described here. (There is the firstgid
	// attribute missing and this source attribute is also not there. These two attributes are kept in the
	// TMX map, since they are map specific.)
	Source string `xml:"source,attr"`
	// External TSX source loaded.
	SourceLoaded bool `xml:"-"`
	// The name of this tileset.
	Name string `xml:"name,attr"`
	// The (maximum) width of the tiles in this tileset.
	TileWidth int `xml:"tilewidth,attr"`
	// The (maximum) height of the tiles in this tileset.
	TileHeight int `xml:"tileheight,attr"`
	// The spacing in pixels between the tiles in this tileset (applies to the tileset image).
	Spacing int `xml:"spacing,attr"`
	// The margin around the tiles in this tileset (applies to the tileset image).
	Margin int `xml:"margin,attr"`
	// The number of tiles in this tileset (since 0.13)
	TileCount int `xml:"tilecount,attr"`
	// The number of tile columns in the tileset. For image collection tilesets it is editable and is used when displaying the tileset. (since 0.15)
	Columns int `xml:"columns,attr"`
	// Offset in pixels, to be applied when drawing a tile from the related tileset. When not present, no offset is applied.
	TileOffset *TilesetTileOffset `xml:"tileoffset"`
	// Custom properties
	Properties Properties `xml:"properties>property"`
	// Embedded image
	Image *Image `xml:"image"`
	// Defines an array of terrain types, which can be referenced from the terrain of the tile element.
	TerrainTypes []*Terrain `xml:"terraintypes>terrain"`
	// Tiles in tileset
	Tiles []*TilesetTile `xml:"tile"`
}

// GetFileFullPath returns path to file relative to tileset file
func (ts *Tileset) GetFileFullPath(fileName string) string {
	return filepath.Join(ts.baseDir, fileName)
}

// TilesetTileOffset is used to specify an offset in pixels, to be applied when drawing a tile from the related tileset. When not present, no offset is applied
type TilesetTileOffset struct {
	// Horizontal offset in pixels
	X int `xml:"x,attr"`
	// Vertical offset in pixels (positive is down)
	Y int `xml:"y,attr"`
}

// Terrain type
type Terrain struct {
	// The name of the terrain type.
	Name string `xml:"name,attr"`
	// The local tile-id of the tile that represents the terrain visually.
	Tile uint32 `xml:"tile,attr"`
	// Custom properties
	Properties Properties `xml:"properties>property"`
}

// TilesetTile information
type TilesetTile struct {
	// The local tile ID within its tileset.
	ID uint32 `xml:"id,attr"`
	// The type of the tile. Refers to an object type and is used by tile objects. (optional) (since 1.0)
	Type string `xml:"type,attr"`
	// Defines the terrain type of each corner of the tile, given as comma-separated indexes in the terrain types
	// array in the order top-left, top-right, bottom-left, bottom-right.
	// Leaving out a value means that corner has no terrain. (optional) (since 0.9)
	Terrain string `xml:"terrain,attr"`
	// A percentage indicating the probability that this tile is chosen when it competes with others while editing with the terrain tool. (optional) (since 0.9)
	Probability int `xml:"probability,attr"`
	// Custom properties
	Properties Properties `xml:"properties>property"`
	// Embedded image
	Image *Image `xml:"image"`
	// Tile object groups
	ObjectGroups []*ObjectGroup `xml:"objectgroup"`
	// List of animation frames
	Animation []*AnimationFrame `xml:"animation>frame"`
}

// AnimationFrame is single frame of animation
type AnimationFrame struct {
	// The local ID of a tile within the parent tileset.
	TileID uint32 `xml:"tileid,attr"`
	// How long (in milliseconds) this frame should be displayed before advancing to the next frame.
	Duration uint32 `xml:"duration,attr"`
}
