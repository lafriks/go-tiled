/*
Copyright (c) 2023 Lauris Bukšis-Haberkorns <lauris@nix.lv>
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

	renderer.RenderObjectGroup(0)

	w, _ := os.Create("../assets/test_render_objects.png")
	defer w.Close()

	if err = renderer.SaveAsPng(w); err != nil {
		t.Error(err)
	}
}
