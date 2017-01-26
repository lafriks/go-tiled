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

package render

import (
	"image"

	"github.com/disintegration/imaging"
	tiled "github.com/lafriks/go-tiled"
)

type OrthogonalRendererEngine struct {
	m *tiled.Map
}

func (e *OrthogonalRendererEngine) Init(m *tiled.Map) {
	e.m = m
}

func (e *OrthogonalRendererEngine) GetFinalImageSize() image.Rectangle {
	return image.Rect(0, 0, e.m.Width*e.m.TileWidth, e.m.Height*e.m.TileHeight)
}

func (e *OrthogonalRendererEngine) RotateTileImage(tile *tiled.LayerTile, img image.Image) image.Image {
	timg := img
	if tile.HorizontalFlip {
		timg = imaging.FlipH(timg)
	}
	if tile.VerticalFlip {
		timg = imaging.FlipV(timg)
	}
	if tile.DiagonalFlip {
		timg = imaging.FlipH(imaging.Rotate270(timg))
	}

	return timg
}

func (e *OrthogonalRendererEngine) GetTilePosition(x, y int) image.Rectangle {
	return image.Rect(x*e.m.TileWidth,
		y*e.m.TileHeight,
		(x+1)*e.m.TileWidth,
		(y+1)*e.m.TileHeight)
}
