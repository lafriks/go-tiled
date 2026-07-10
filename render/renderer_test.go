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
	"os"
	"path/filepath"
	"testing"

	"github.com/lafriks/go-tiled"
)

func TestRenderer_RenderOrthogonalMap(t *testing.T) {
	tiledMap, err := tiled.LoadFile("../assets/test_wangsets_map.tmx")
	if err != nil {
		t.Error(err)
		return
	}

	renderer, err := NewRenderer(tiledMap)
	if err != nil {
		t.Error(err)
		return
	}

	if err = renderer.RenderVisibleLayers(); err != nil {
		t.Error(err)
		return
	}

	w, err := os.Create(filepath.Join(t.TempDir(), "test_render_orthogonal.png"))
	if err != nil {
		t.Error(err)
		return
	}
	defer w.Close()

	if err = renderer.SaveAsPng(w); err != nil {
		t.Error(err)
	}
}

func TestRenderer_RenderIsometricMap(t *testing.T) {
	tiledMap, err := tiled.LoadFile("../assets/test_isometric.tmx")
	if err != nil {
		t.Error(err)
		return
	}

	renderer, err := NewRenderer(tiledMap)
	if err != nil {
		t.Error(err)
		return
	}

	if err = renderer.RenderVisibleLayers(); err != nil {
		t.Error(err)
		return
	}

	if got := renderer.Result.Bounds().Dx(); got != 800 {
		t.Errorf("image width = %d, want 800", got)
	}
	if got := renderer.Result.Bounds().Dy(); got != 400 {
		t.Errorf("image height = %d, want 400", got)
	}

	w, err := os.Create(filepath.Join(t.TempDir(), "test_render_isometric.png"))
	if err != nil {
		t.Error(err)
		return
	}
	defer w.Close()

	if err = renderer.SaveAsPng(w); err != nil {
		t.Error(err)
	}
}
