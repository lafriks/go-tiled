package render

import (
	"github.com/lafriks/go-tiled"
	"os"
	"testing"
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
