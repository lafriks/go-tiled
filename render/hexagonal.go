/*
Copyright (c) 2026 Lauris Bukšis <lauris@nix.lv>

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

package render

import (
	"image"

	"github.com/disintegration/imaging"
	"github.com/lafriks/go-tiled"
)

// HexagonalRendererEngine represents a hexagonal rendering engine (Tiled's
// "Hexagonal (Staggered)" map type, orientation="hexagonal"), supporting both
// pointy-top (StaggerAxis "y") and flat-top (StaggerAxis "x") variants, with
// an arbitrary HexSideLength and either StaggerIndex. See StaggeredRendererEngine
// for orientation="staggered" (Tiled's "Isometric (Staggered)" map type),
// which is not a hexagonal map but shares this same geometry engine.
//
// The geometry below mirrors Tiled's own HexagonalRenderer::RenderParams and
// tileToScreenCoords/boundingRect exactly (see
// https://github.com/mapeditor/tiled/blob/master/src/libtiled/hexagonalrenderer.cpp),
// rather than assuming HexSideLength is always exactly half of TileWidth/TileHeight.
type HexagonalRendererEngine struct {
	m *tiled.Map

	staggerX    bool
	staggerEven bool

	// sideOffsetX/Y is the flat "cap" on each side of the hex that isn't part of
	// the staggered side length; columnWidth/rowHeight is the step between
	// adjacent columns/rows along the staggered axis.
	sideOffsetX, sideOffsetY int
	columnWidth, rowHeight   int
}

// Init initializes rendering engine with provided map options.
func (e *HexagonalRendererEngine) Init(m *tiled.Map) {
	e.m = m
	e.staggerX = m.StaggerAxis == tiled.AxisX
	e.staggerEven = m.StaggerIndex == tiled.StaggerIndexEven

	sideLengthX, sideLengthY := 0, 0
	if e.staggerX {
		sideLengthX = m.HexSideLength
	} else {
		sideLengthY = m.HexSideLength
	}

	e.sideOffsetX = (m.TileWidth - sideLengthX) / 2
	e.sideOffsetY = (m.TileHeight - sideLengthY) / 2
	e.columnWidth = e.sideOffsetX + sideLengthX
	e.rowHeight = e.sideOffsetY + sideLengthY
}

// doStagger reports whether the given column (for StaggerAxis "x") or row (for
// StaggerAxis "y") index is one of the ones shifted half a step, per StaggerIndex.
func (e *HexagonalRendererEngine) doStagger(index int) bool {
	return (index&1 == 1) != e.staggerEven
}

// GetFinalImageSize returns final image size based on map data.
func (e *HexagonalRendererEngine) GetFinalImageSize() image.Rectangle {
	var width, height int

	if e.staggerX {
		width = e.m.Width*e.columnWidth + e.sideOffsetX
		height = e.m.Height * e.m.TileHeight
		if e.m.Width > 1 {
			height += e.rowHeight
		}
	} else {
		width = e.m.Width * e.m.TileWidth
		height = e.m.Height*e.rowHeight + e.sideOffsetY
		if e.m.Height > 1 {
			width += e.columnWidth
		}
	}

	return image.Rect(0, 0, width, height)
}

// RotateTileImage rotates provided tile layer.
func (e *HexagonalRendererEngine) RotateTileImage(tile *tiled.LayerTile, img image.Image) image.Image {
	timg := img
	if tile.DiagonalFlip {
		timg = imaging.FlipH(imaging.Rotate270(timg))
	}
	if tile.HorizontalFlip {
		timg = imaging.FlipH(timg)
	}
	if tile.VerticalFlip {
		timg = imaging.FlipV(timg)
	}

	return timg
}

// GetTilePosition returns tile position in image.
//
// As with the orthogonal/isometric engines, a tile image taller than the map's
// TileHeight is anchored to the bottom of its (TileHeight-tall) grid cell per
// the Tiled spec, using imgSize rather than assuming it matches TileHeight.
func (e *HexagonalRendererEngine) GetTilePosition(x, y int, imgSize image.Point) image.Point {
	var px, py int

	if e.staggerX {
		px = x * e.columnWidth
		py = y * e.m.TileHeight
		if e.doStagger(x) {
			py += e.rowHeight
		}
	} else {
		px = x * e.m.TileWidth
		if e.doStagger(y) {
			px += e.columnWidth
		}
		py = y * e.rowHeight
	}

	return image.Pt(px, py+e.m.TileHeight-imgSize.Y)
}

// PixelToScreenCoords returns screen coordinates for a raw object pixel position.
// Tiled's HexagonalRenderer doesn't project object positions (unlike isometric);
// they're stored directly in screen pixel space, so this is the identity transform.
func (e *HexagonalRendererEngine) PixelToScreenCoords(x, y float64) (float64, float64) {
	return x, y
}

// GetObjectAnchor returns the bottom-left point of a tile object's image. Tiled's
// default tile-object alignment is bottom-center only for isometric maps; every
// other orientation, including hexagonal, defaults to bottom-left.
func (e *HexagonalRendererEngine) GetObjectAnchor(imgSize image.Point) image.Point {
	return image.Pt(0, imgSize.Y-1)
}
