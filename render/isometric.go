package render

import (
	"image"

	"github.com/disintegration/imaging"
	"github.com/lafriks/go-tiled"
)

type IsometricRendererEngine struct {
	m *tiled.Map
}

func (e *IsometricRendererEngine) Init(m *tiled.Map) {
	e.m = m
}

func (e *IsometricRendererEngine) GetFinalImageSize() image.Rectangle {
	return image.Rect(0, 0, e.m.Width*e.m.TileWidth, e.m.Height*e.m.TileHeight)
}

func (e *IsometricRendererEngine) RotateTileImage(tile *tiled.LayerTile, img image.Image) image.Image {
	timg := img
	if tile.DiagonalFlip {
		timg = imaging.FlipH(imaging.Rotate270(timg))
	}
	if tile.HorizontalFlip {
		timg = imaging.FlipH(timg)
	}
	if tile.VerticalFlip {
		timg = imaging.FlipV(timg)
	}

	return timg
}

func (e *IsometricRendererEngine) GetTilePosition(x, y int) image.Rectangle {
    tw, th := e.m.TileWidth, e.m.TileHeight

	ratio := tw / th

	stepX := tw / ratio
	stepY := th / ratio

	actualSpriteHeight := th * ratio

    offsetX := e.m.Width * stepX
	offsetY := -th

    sx := (x - y) * stepX + offsetX - (tw / ratio)
    sy := (x + y) * stepY + offsetY

    return image.Rect(sx, sy, sx + tw, sy + actualSpriteHeight)
}