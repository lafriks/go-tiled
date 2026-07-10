/*
Copyright (c) 2017 Lauris Bukšis <lauris@nix.lv>

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
	tiled "github.com/lafriks/go-tiled"
)

// OrthogonalRendererEngine represents orthogonal rendering engine.
type OrthogonalRendererEngine struct {
	m *tiled.Map
}

// Init initializes rendering engine with provided map options.
func (e *OrthogonalRendererEngine) Init(m *tiled.Map) {
	e.m = m
}

// GetFinalImageSize returns final image size based on map data.
func (e *OrthogonalRendererEngine) GetFinalImageSize() image.Rectangle {
	return image.Rect(0, 0, e.m.Width*e.m.TileWidth, e.m.Height*e.m.TileHeight)
}

// RotateTileImage rotates provided tile layer.
func (e *OrthogonalRendererEngine) RotateTileImage(tile *tiled.LayerTile, img image.Image) image.Image {
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
// Per the Tiled spec, a tile image larger than the map's grid cell "extends
// at the top and right (anchored to the bottom left)", so the image's bottom
// edge is aligned to the bottom of the grid cell rather than its top.
func (e *OrthogonalRendererEngine) GetTilePosition(x, y int, imgSize image.Point) image.Point {
	return image.Pt(x*e.m.TileWidth, (y+1)*e.m.TileHeight-imgSize.Y)
}

// PixelToScreenCoords returns screen coordinates for a raw object pixel position.
// Orthogonal maps store object positions directly in screen pixel space already,
// so this is the identity transform.
func (e *OrthogonalRendererEngine) PixelToScreenCoords(x, y float64) (float64, float64) {
	return x, y
}

// GetObjectAnchor returns the bottom-left point of a tile object's image, per the
// Tiled spec's orthogonal alignment rule.
func (e *OrthogonalRendererEngine) GetObjectAnchor(imgSize image.Point) image.Point {
	return image.Pt(0, imgSize.Y-1)
}
