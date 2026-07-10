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
	tiled "github.com/lafriks/go-tiled"
)

// IsometricRendererEngine represents isometric rendering engine.
type IsometricRendererEngine struct {
	m *tiled.Map
}

// Init initializes rendering engine with provided map options.
func (e *IsometricRendererEngine) Init(m *tiled.Map) {
	e.m = m
}

// GetFinalImageSize returns final image size based on map data.
func (e *IsometricRendererEngine) GetFinalImageSize() image.Rectangle {
	side := e.m.Height + e.m.Width
	hx := side * e.m.TileWidth / 2
	hy := side * e.m.TileHeight / 2

	return image.Rect(0, 0, hx, hy)
}

// RotateTileImage rotates provided tile layer.
func (e *IsometricRendererEngine) RotateTileImage(tile *tiled.LayerTile, img image.Image) image.Image {
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
// The horizontal component matches Tiled's own IsometricRenderer::tileToScreenCoords
// (shifted left by half a tile width so the image is centered on the diamond grid
// vertex, as Tiled's drawTileLayer does). The vertical component anchors the image to
// the bottom of the tile's diamond footprint: imgSize.Y is used rather than assuming
// it equals the map's TileHeight, so tilesets mixing flat ground tiles with taller
// wall/structure tiles (of any height) are positioned correctly.
func (e *IsometricRendererEngine) GetTilePosition(x, y int, imgSize image.Point) image.Point {
	tw, th := e.m.TileWidth, e.m.TileHeight

	stepX := tw / 2
	stepY := th / 2

	offsetX := e.m.Height * tw / 2

	sx := (x-y)*stepX + offsetX - stepX
	sy := (x+y)*stepY + th - imgSize.Y

	return image.Pt(sx, sy)
}

// PixelToScreenCoords converts a raw object pixel position (as stored in the TMX
// file) into final screen coordinates, matching Tiled's own
// IsometricRenderer::pixelToScreenCoords exactly.
//
// Note that x is divided by TileHeight, not TileWidth: Tiled treats an isometric
// map's unprojected object space as using a square unit cell sized by TileHeight
// for both axes, and only brings in TileWidth for the final horizontal scaling.
// This is confirmed by Tiled's own source and its maintainer (see
// https://discourse.mapeditor.org/t/whats-the-algorithm-of-object-position-in-iso-map/1790).
func (e *IsometricRendererEngine) PixelToScreenCoords(x, y float64) (float64, float64) {
	tw := float64(e.m.TileWidth)
	th := float64(e.m.TileHeight)
	originX := float64(e.m.Height) * tw / 2

	tileX := x / th
	tileY := y / th

	sx := (tileX-tileY)*tw/2 + originX
	sy := (tileX + tileY) * th / 2

	return sx, sy
}

// GetObjectAnchor returns the bottom-center point of a tile object's image, per
// the Tiled spec's isometric alignment rule.
func (e *IsometricRendererEngine) GetObjectAnchor(imgSize image.Point) image.Point {
	return image.Pt(imgSize.X/2, imgSize.Y-1)
}
