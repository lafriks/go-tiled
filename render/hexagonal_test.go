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
	"testing"

	"github.com/lafriks/go-tiled"
)

// TestHexagonalRendererEngine_GetTilePosition_FlatTop pins down flat-top
// (StaggerAxis "x") tile placement against Tiled's own
// HexagonalRenderer::tileToScreenCoords/RenderParams. TileWidth, TileHeight and
// HexSideLength are chosen so that HexSideLength is NOT half of TileWidth (as it
// happens to be in assets/hex.tmx), so a formula that assumes that ratio (as a
// naive "3/4 of tile size" approximation would) fails this test instead of only
// coincidentally passing.
func TestHexagonalRendererEngine_GetTilePosition_FlatTop(t *testing.T) {
	m := &tiled.Map{
		Width: 6, Height: 5, TileWidth: 40, TileHeight: 20, HexSideLength: 8,
		StaggerAxis: tiled.AxisX, StaggerIndex: tiled.StaggerIndexOdd,
	}
	e := &HexagonalRendererEngine{}
	e.Init(m)

	tests := []struct {
		name    string
		x, y    int
		imgSize image.Point
		want    image.Point
	}{
		{"(0,0) even column, flat image", 0, 0, image.Pt(40, 20), image.Pt(0, 0)},
		{"(1,0) odd column (staggered), flat image", 1, 0, image.Pt(40, 20), image.Pt(24, 10)},
		{"(2,0) even column, flat image", 2, 0, image.Pt(40, 20), image.Pt(48, 0)},
		{"(1,3) odd column, flat image", 1, 3, image.Pt(40, 20), image.Pt(24, 70)},
		{"(1,0) odd column, tall image", 1, 0, image.Pt(40, 30), image.Pt(24, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.GetTilePosition(tt.x, tt.y, tt.imgSize); got != tt.want {
				t.Errorf("GetTilePosition(%d, %d, %v) = %v, want %v", tt.x, tt.y, tt.imgSize, got, tt.want)
			}
		})
	}

	if got, want := e.GetFinalImageSize(), image.Rect(0, 0, 160, 110); got != want {
		t.Errorf("GetFinalImageSize() = %v, want %v", got, want)
	}
}

// TestHexagonalRendererEngine_GetTilePosition_PointyTop mirrors the flat-top
// test above for StaggerAxis "y", additionally using StaggerIndex "even" (vs
// "odd" above) to confirm StaggerIndex is honored on both axes.
func TestHexagonalRendererEngine_GetTilePosition_PointyTop(t *testing.T) {
	m := &tiled.Map{
		Width: 5, Height: 6, TileWidth: 40, TileHeight: 30, HexSideLength: 10,
		StaggerAxis: tiled.AxisY, StaggerIndex: tiled.StaggerIndexEven,
	}
	e := &HexagonalRendererEngine{}
	e.Init(m)

	tests := []struct {
		name    string
		x, y    int
		imgSize image.Point
		want    image.Point
	}{
		{"(0,0) even row (staggered), flat image", 0, 0, image.Pt(40, 30), image.Pt(20, 0)},
		{"(0,1) odd row, flat image", 0, 1, image.Pt(40, 30), image.Pt(0, 20)},
		{"(0,2) even row (staggered), flat image", 0, 2, image.Pt(40, 30), image.Pt(20, 40)},
		{"(3,1) odd row, flat image", 3, 1, image.Pt(40, 30), image.Pt(120, 20)},
		{"(0,1) odd row, tall image", 0, 1, image.Pt(40, 50), image.Pt(0, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.GetTilePosition(tt.x, tt.y, tt.imgSize); got != tt.want {
				t.Errorf("GetTilePosition(%d, %d, %v) = %v, want %v", tt.x, tt.y, tt.imgSize, got, tt.want)
			}
		})
	}

	if got, want := e.GetFinalImageSize(), image.Rect(0, 0, 220, 130); got != want {
		t.Errorf("GetFinalImageSize() = %v, want %v", got, want)
	}
}

// TestHexagonalRendererEngine_PixelToScreenCoords verifies object positions
// pass through unchanged: Tiled's HexagonalRenderer doesn't override
// pixelToScreenCoords, so hex/staggered maps store object positions directly
// in screen pixel space, same as orthogonal.
func TestHexagonalRendererEngine_PixelToScreenCoords(t *testing.T) {
	m := &tiled.Map{Width: 10, Height: 10, TileWidth: 32, TileHeight: 32, HexSideLength: 16, StaggerAxis: tiled.AxisY}
	e := &HexagonalRendererEngine{}
	e.Init(m)

	gx, gy := e.PixelToScreenCoords(123, 456)
	if gx != 123 || gy != 456 {
		t.Errorf("PixelToScreenCoords(123, 456) = (%v, %v), want (123, 456)", gx, gy)
	}
}

// TestHexagonalRendererEngine_GetObjectAnchor verifies hexagonal tile objects
// default to bottom-left alignment, per Tiled's MapObject::alignment (which
// only special-cases Isometric orientation for bottom-center).
func TestHexagonalRendererEngine_GetObjectAnchor(t *testing.T) {
	e := &HexagonalRendererEngine{}
	e.Init(&tiled.Map{TileWidth: 32, TileHeight: 32, HexSideLength: 16, StaggerAxis: tiled.AxisY})

	if got, want := e.GetObjectAnchor(image.Pt(10, 6)), image.Pt(0, 5); got != want {
		t.Errorf("GetObjectAnchor(10, 6) = %v, want %v", got, want)
	}
}
