package render

import (
	"os"
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

	renderer.RenderVisibleLayers()

	w, _ := os.Create("../assets/test_orthogonal.png")
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

	renderer.RenderVisibleLayers()
	//renderer.RenderLayer(0)

	w, _ := os.Create("../assets/test_render_isomap.png")
	defer w.Close()

	if err = renderer.SaveAsPng(w); err != nil {
		t.Error(err)
	}

}