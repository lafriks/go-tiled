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

// TestStaggeredRendererEngine_GetTilePosition_StaggerX pins down StaggerAxis
// "x" tile placement. A staggered map's HexSideLength is always 0 (the
// attribute doesn't apply to this orientation), which reduces
// HexagonalRendererEngine's geometry to a plain half-tile brick offset; these
// numbers were independently derived and cross-checked against that formula.
func TestStaggeredRendererEngine_GetTilePosition_StaggerX(t *testing.T) {
	m := &tiled.Map{
		Width: 6, Height: 5, TileWidth: 40, TileHeight: 20,
		StaggerAxis: tiled.AxisX, StaggerIndex: tiled.StaggerIndexOdd,
	}
	e := &StaggeredRendererEngine{}
	e.Init(m)

	tests := []struct {
		name    string
		x, y    int
		imgSize image.Point
		want    image.Point
	}{
		{"(0,0) even column, flat image", 0, 0, image.Pt(40, 20), image.Pt(0, 0)},
		{"(1,0) odd column (staggered), flat image", 1, 0, image.Pt(40, 20), image.Pt(20, 10)},
		{"(2,0) even column, flat image", 2, 0, image.Pt(40, 20), image.Pt(40, 0)},
		{"(1,3) odd column, flat image", 1, 3, image.Pt(40, 20), image.Pt(20, 70)},
		{"(1,0) odd column, tall image", 1, 0, image.Pt(40, 30), image.Pt(20, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.GetTilePosition(tt.x, tt.y, tt.imgSize); got != tt.want {
				t.Errorf("GetTilePosition(%d, %d, %v) = %v, want %v", tt.x, tt.y, tt.imgSize, got, tt.want)
			}
		})
	}

	if got, want := e.GetFinalImageSize(), image.Rect(0, 0, 140, 110); got != want {
		t.Errorf("GetFinalImageSize() = %v, want %v", got, want)
	}
}

// TestStaggeredRendererEngine_GetTilePosition_StaggerY mirrors the StaggerX
// test above for StaggerAxis "y", using StaggerIndex "even" to confirm
// StaggerIndex is honored on both axes.
func TestStaggeredRendererEngine_GetTilePosition_StaggerY(t *testing.T) {
	m := &tiled.Map{
		Width: 5, Height: 6, TileWidth: 40, TileHeight: 30,
		StaggerAxis: tiled.AxisY, StaggerIndex: tiled.StaggerIndexEven,
	}
	e := &StaggeredRendererEngine{}
	e.Init(m)

	tests := []struct {
		name    string
		x, y    int
		imgSize image.Point
		want    image.Point
	}{
		{"(0,0) even row (staggered), flat image", 0, 0, image.Pt(40, 30), image.Pt(20, 0)},
		{"(0,1) odd row, flat image", 0, 1, image.Pt(40, 30), image.Pt(0, 15)},
		{"(0,2) even row (staggered), flat image", 0, 2, image.Pt(40, 30), image.Pt(20, 30)},
		{"(3,1) odd row, flat image", 3, 1, image.Pt(40, 30), image.Pt(120, 15)},
		{"(0,1) odd row, tall image", 0, 1, image.Pt(40, 50), image.Pt(0, -5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.GetTilePosition(tt.x, tt.y, tt.imgSize); got != tt.want {
				t.Errorf("GetTilePosition(%d, %d, %v) = %v, want %v", tt.x, tt.y, tt.imgSize, got, tt.want)
			}
		})
	}

	if got, want := e.GetFinalImageSize(), image.Rect(0, 0, 220, 105); got != want {
		t.Errorf("GetFinalImageSize() = %v, want %v", got, want)
	}
}

// TestStaggeredRendererEngine_PixelToScreenCoords and
// TestStaggeredRendererEngine_GetObjectAnchor exercise the methods
// StaggeredRendererEngine gets for free via embedding HexagonalRendererEngine,
// confirming method promotion actually wires them up as expected.
func TestStaggeredRendererEngine_PixelToScreenCoords(t *testing.T) {
	e := &StaggeredRendererEngine{}
	e.Init(&tiled.Map{Width: 10, Height: 10, TileWidth: 32, TileHeight: 32, StaggerAxis: tiled.AxisY})

	gx, gy := e.PixelToScreenCoords(123, 456)
	if gx != 123 || gy != 456 {
		t.Errorf("PixelToScreenCoords(123, 456) = (%v, %v), want (123, 456)", gx, gy)
	}
}

func TestStaggeredRendererEngine_GetObjectAnchor(t *testing.T) {
	e := &StaggeredRendererEngine{}
	e.Init(&tiled.Map{TileWidth: 32, TileHeight: 32, StaggerAxis: tiled.AxisY})

	if got, want := e.GetObjectAnchor(image.Pt(10, 6)), image.Pt(0, 5); got != want {
		t.Errorf("GetObjectAnchor(10, 6) = %v, want %v", got, want)
	}
}
