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
	"testing"

	"github.com/lafriks/go-tiled"
)

// TestIsometricRendererEngine_GetTilePosition pins down the tile placement
// formula against Tiled's own IsometricRenderer reference implementation:
//
//	screenY = (x+y)*tileHeight/2 + tileHeight - imageHeight
//	screenX = (x-y)*tileWidth/2 + mapHeight*tileWidth/2 - tileWidth/2
//
// TileWidth (64) and TileHeight (32) are deliberately unequal, and the "tall"
// case uses an image height that differs from both, so a formula that
// substitutes TileWidth for the image height (as the previous, buggy version
// did) would fail this test instead of only coincidentally passing.
func TestIsometricRendererEngine_GetTilePosition(t *testing.T) {
	m := &tiled.Map{Width: 10, Height: 8, TileWidth: 64, TileHeight: 32}
	e := &IsometricRendererEngine{}
	e.Init(m)

	tests := []struct {
		name    string
		x, y    int
		imgSize image.Point
		want    image.Point
	}{
		{"origin, flat tile (imgHeight == tileHeight)", 0, 0, image.Pt(64, 32), image.Pt(224, 0)},
		{"origin, tall tile (imgHeight == 48)", 0, 0, image.Pt(64, 48), image.Pt(224, -16)},
		{"(3,2), flat tile", 3, 2, image.Pt(64, 32), image.Pt(256, 80)},
		{"(3,2), tall tile", 3, 2, image.Pt(64, 48), image.Pt(256, 64)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.GetTilePosition(tt.x, tt.y, tt.imgSize); got != tt.want {
				t.Errorf("GetTilePosition(%d, %d, %v) = %v, want %v", tt.x, tt.y, tt.imgSize, got, tt.want)
			}
		})
	}
}

// TestIsometricRendererEngine_PixelToScreenCoords pins down object-position
// projection against Tiled's own IsometricRenderer::pixelToScreenCoords, as
// confirmed by Tiled's maintainer:
// https://discourse.mapeditor.org/t/whats-the-algorithm-of-object-position-in-iso-map/1790
//
//	tileX = x / tileHeight   (not tileWidth -- both axes use a tileHeight-sized unit)
//	tileY = y / tileHeight
//	screenX = (tileX-tileY)*tileWidth/2 + mapHeight*tileWidth/2
//	screenY = (tileX+tileY)*tileHeight/2
//
// TileWidth (64) and TileHeight (32) are deliberately unequal so that dividing
// x by the wrong dimension would produce a visibly different (and wrong) result.
func TestIsometricRendererEngine_PixelToScreenCoords(t *testing.T) {
	m := &tiled.Map{Width: 10, Height: 8, TileWidth: 64, TileHeight: 32}
	e := &IsometricRendererEngine{}
	e.Init(m)

	tests := []struct {
		name string
		x, y float64
		wx   float64
		wy   float64
	}{
		{"origin", 0, 0, 256, 0},
		{"(128,64)", 128, 64, 320, 96},
		{"(96,32)", 96, 32, 320, 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gx, gy := e.PixelToScreenCoords(tt.x, tt.y)
			if gx != tt.wx || gy != tt.wy {
				t.Errorf("PixelToScreenCoords(%v, %v) = (%v, %v), want (%v, %v)", tt.x, tt.y, gx, gy, tt.wx, tt.wy)
			}
		})
	}
}

// TestOrthogonalRendererEngine_PixelToScreenCoords verifies orthogonal object
// positions pass through unchanged, since they're already stored in screen
// pixel space per the Tiled spec.
func TestOrthogonalRendererEngine_PixelToScreenCoords(t *testing.T) {
	m := &tiled.Map{Width: 10, Height: 8, TileWidth: 64, TileHeight: 32}
	e := &OrthogonalRendererEngine{}
	e.Init(m)

	gx, gy := e.PixelToScreenCoords(123, 456)
	if gx != 123 || gy != 456 {
		t.Errorf("PixelToScreenCoords(123, 456) = (%v, %v), want (123, 456)", gx, gy)
	}
}

// TestGetObjectAnchor verifies each engine picks the tile-object alignment
// point mandated by the Tiled spec: bottom-left for orthogonal, bottom-center
// for isometric.
func TestGetObjectAnchor(t *testing.T) {
	imgSize := image.Pt(10, 6)

	ortho := &OrthogonalRendererEngine{}
	ortho.Init(&tiled.Map{})
	if got, want := ortho.GetObjectAnchor(imgSize), image.Pt(0, 5); got != want {
		t.Errorf("OrthogonalRendererEngine.GetObjectAnchor(%v) = %v, want %v", imgSize, got, want)
	}

	iso := &IsometricRendererEngine{}
	iso.Init(&tiled.Map{})
	if got, want := iso.GetObjectAnchor(imgSize), image.Pt(5, 5); got != want {
		t.Errorf("IsometricRendererEngine.GetObjectAnchor(%v) = %v, want %v", imgSize, got, want)
	}
}
