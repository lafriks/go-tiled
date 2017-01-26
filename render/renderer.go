/*
Copyright (c) 2017 Lauris Bukšis-Haberkorns <lauris@nix.lv>

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
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"

	"image/jpeg"

	"image/gif"

	"github.com/disintegration/imaging"
	"github.com/lafriks/go-tiled"
)

var (
	ErrUnsupportedOrientation = errors.New("tiled/render: unsupported orientation")
	ErrUnsupportedRenderOrder = errors.New("tiled/render: unsupported render order")
)

type RendererEngine interface {
	Init(m *tiled.Map)
	GetFinalImageSize() image.Rectangle
	RotateTileImage(tile *tiled.LayerTile, img image.Image) image.Image
	GetTilePosition(x, y int) image.Rectangle
}

type Renderer struct {
	m         *tiled.Map
	img       *image.NRGBA
	tileCache map[uint32]image.Image
	engine    RendererEngine
}

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

func NewRenderer(m *tiled.Map) (*Renderer, error) {
	r := &Renderer{m: m, tileCache: make(map[uint32]image.Image)}
	if r.m.Orientation == "orthogonal" {
		r.engine = &OrthogonalRendererEngine{}
	} else {
		return nil, ErrUnsupportedOrientation
	}

	r.engine.Init(r.m)
	r.img = image.NewNRGBA(r.engine.GetFinalImageSize())

	return r, nil
}

func (r *Renderer) getTileImage(tile *tiled.LayerTile) (image.Image, error) {
	timg, ok := r.tileCache[tile.Tileset.FirstGID+tile.ID]
	// Precache all tiles in tileset
	if !ok {
		sf, err := os.Open(r.m.GetFileFullPath(tile.Tileset.Image.Source))
		if err != nil {
			return nil, err
		}
		defer sf.Close()

		img, _, err := image.Decode(sf)
		if err != nil {
			return nil, err
		}

		tilesetTileCount := tile.Tileset.TileCount

		tilesetColumns := tile.Tileset.Columns
		if tilesetColumns == 0 {
			tilesetColumns = tile.Tileset.Image.Width / (tile.Tileset.TileWidth + tile.Tileset.Spacing)
		}

		if tilesetTileCount == 0 {
			tilesetTileCount = (tile.Tileset.Image.Height / (tile.Tileset.TileHeight + tile.Tileset.Spacing)) * tilesetColumns
		}

		for i := tile.Tileset.FirstGID; i < tile.Tileset.FirstGID+uint32(tilesetTileCount); i++ {
			x := int(i-tile.Tileset.FirstGID) % tilesetColumns
			y := int(i-tile.Tileset.FirstGID) / tilesetColumns

			rect := image.Rect(x*tile.Tileset.TileWidth,
				y*tile.Tileset.TileHeight,
				(x+1)*tile.Tileset.TileWidth,
				(y+1)*tile.Tileset.TileHeight)

			r.tileCache[i] = imaging.Crop(img, rect)
			if tile.ID == i-tile.Tileset.FirstGID {
				timg = r.tileCache[i]
			}
		}
	}

	timg = r.engine.RotateTileImage(tile, timg)

	return timg, nil
}

func (r *Renderer) RenderLayer(index int) error {
	layer := r.m.Layers[index]

	var xs, xe, xi, ys, ye, yi int
	if r.m.RenderOrder == "" || r.m.RenderOrder == "right-down" {
		xs = 0
		xe = r.m.Width
		xi = 1
		ys = 0
		ye = r.m.Height
		yi = 1
	} else {
		return ErrUnsupportedRenderOrder
	}

	i := 0
	for y := ys; y*yi < ye; y = y + yi {
		for x := xs; x*xi < xe; x = x + xi {
			if layer.Tiles[i].IsNil() {
				i++
				continue
			}

			img, err := r.getTileImage(layer.Tiles[i])
			if err != nil {
				return err
			}

			pos := r.engine.GetTilePosition(x, y)

			if layer.Opacity < 1 {
				mask := image.NewUniform(color.Alpha{uint8(layer.Opacity * 255)})

				draw.DrawMask(r.img, pos, img, img.Bounds().Min, mask, mask.Bounds().Min, draw.Over)
			} else {
				draw.Draw(r.img, pos, img, img.Bounds().Min, draw.Over)
			}

			i++
		}
	}

	return nil
}

func (r *Renderer) RenderVisibleLayers() error {
	for i := range r.m.Layers {
		if !r.m.Layers[i].Visible {
			continue
		}

		if err := r.RenderLayer(i); err != nil {
			return err
		}
	}

	return nil
}

func (r *Renderer) SaveAsPng(w io.Writer) error {
	return png.Encode(w, r.img)
}

func (r *Renderer) SaveAsJpeg(w io.Writer, options *jpeg.Options) error {
	return jpeg.Encode(w, r.img, options)
}

func (r *Renderer) SaveAsGif(w io.Writer, options *gif.Options) error {
	return gif.Encode(w, r.img, options)
}
