package render

import (
	"image"
	"os"
	"testing"
	"path/filepath"

	"github.com/lafriks/go-tiled"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	dir := filepath.Join("..", "assets", "test_output")
	if err := os.MkdirAll(dir, 0755); err != nil {
		os.Exit(1)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

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

	w, _ := os.Create("../assets/test_output/test_render_orthogonal.png")
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

	outputPath := "../assets/test_output/test_render_isomap.png"

	w, _ := os.Create(outputPath)
	defer w.Close()

	if err = renderer.SaveAsPng(w); err != nil {
		t.Error(err)
	}

	file, err := os.Open(outputPath)
	require.NoError(t, err)
	defer file.Close()

	img, _, err := image.Decode(file)
	require.NoError(t, err)

	assert.Equal(t, 800, img.Bounds().Dx(), "image width should be 800 pixels")
	assert.Equal(t, 400, img.Bounds().Dy(), "image height should be 400 pixels")
}