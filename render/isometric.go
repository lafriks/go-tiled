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
	side := e.m.Height + e.m.Width
	hx := side * e.m.TileWidth/2
	hy := side * e.m.TileHeight/2

	return image.Rect(0, 0, hx, hy)
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

func (e *IsometricRendererEngine) GetTilePosition(x, y int) image.Point {
    tw, th := e.m.TileWidth, e.m.TileHeight

	stepX := tw / 2
	stepY := th / 2

	//actualSpriteHeight := th * 2

    offsetX := e.m.Height * e.m.TileWidth/2

	offsetY := 0
	if tw > th {
		offsetY = -th
	}

    sx := (x - y) * stepX + offsetX - stepX
    sy := (x + y) * stepY + offsetY

    return image.Pt(sx, sy)
}