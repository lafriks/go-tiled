/*
Copyright (c) 2023 Lauris Bukšis <lauris@nix.lv>
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
	"os"
	"path/filepath"
	"testing"

	"github.com/lafriks/go-tiled"
)

func TestRenderer_RenderObjectGroup(t *testing.T) {
	tiledMap, err := tiled.LoadFile("../assets/test_render_objects.tmx")
	if err != nil {
		t.Error(err)
		return
	}

	renderer, err := NewRenderer(tiledMap)
	if err != nil {
		t.Error(err)
		return
	}

	if err = renderer.RenderObjectGroup(0); err != nil {
		t.Error(err)
		return
	}

	w, err := os.Create(filepath.Join(t.TempDir(), "test_render_objects.png"))
	if err != nil {
		t.Error(err)
		return
	}
	defer w.Close()

	if err = renderer.SaveAsPng(w); err != nil {
		t.Error(err)
	}
}

// TestRenderer_rotateObjectImage_anchor verifies _rotateObjectImage re-anchors
// around whichever point is passed in (rather than always the bottom-left
// corner), which is what lets tile-object rendering support the Tiled spec's
// per-orientation alignment rule (bottom-left for orthogonal, bottom-center
// for isometric).
func TestRenderer_rotateObjectImage_anchor(t *testing.T) {
	r := &Renderer{}
	img := image.NewNRGBA(image.Rect(0, 0, 10, 6))

	tests := []struct {
		name     string
		rotation float64
		anchor   image.Point
		want     image.Point
	}{
		{"no rotation, bottom-left", 0, image.Pt(0, 5), image.Pt(0, 5)},
		{"no rotation, bottom-center", 0, image.Pt(5, 5), image.Pt(5, 5)},
		{"90 degrees, bottom-left", 90, image.Pt(0, 5), image.Pt(0, 0)},
		{"90 degrees, bottom-center", 90, image.Pt(5, 5), image.Pt(0, 5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := r._rotateObjectImage(img, tt.rotation, tt.anchor)
			if got != tt.want {
				t.Errorf("_rotateObjectImage(rotation=%v, anchor=%v) origin = %v, want %v", tt.rotation, tt.anchor, got, tt.want)
			}
		})
	}
}
