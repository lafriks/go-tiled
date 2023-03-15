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
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/internal/utils"

	"github.com/disintegration/imaging"
)

// RenderVisibleGroups renders all visible groups
func (r *Renderer) RenderVisibleGroups() error {
	for _, group := range r.m.Groups {
		if !group.Visible {
			continue
		}
		if err := r._renderGroup(group); err != nil {
			return err
		}
	}
	return nil
}

// RenderGroup renders single group.
func (r *Renderer) RenderGroup(groupIdx int) error {
	if groupIdx >= len(r.m.Groups) {
		return ErrOutOfBounds
	}

	group := r.m.Groups[groupIdx]
	return r._renderGroup(group)
}

func (r *Renderer) _renderGroup(group *tiled.Group) error {
	for _, layer := range group.Layers {
		if !layer.Visible {
			continue
		}
		if err := r._renderLayer(layer); err != nil {
			return err
		}
	}

	for _, objectGroup := range group.ObjectGroups {
		if !objectGroup.Visible {
			continue
		}
		if err := r._renderObjectGroup(objectGroup); err != nil {
			return err
		}
	}

	return nil
}

// RenderVisibleLayersAndObjectGroups render all layers and object groups, layer first, objectGroup second
// so the order may be incorrect,
// you may put them into different groups, then call RenderVisibleGroups
func (r *Renderer) RenderVisibleLayersAndObjectGroups() error {
	// TODO: The order maybe incorrect

	if err := r.RenderVisibleLayers(); err != nil {
		return err
	}
	return r.RenderVisibleObjectGroups()
}

// RenderVisibleObjectGroups renders all visible object groups
func (r *Renderer) RenderVisibleObjectGroups() error {
	for i, layer := range r.m.ObjectGroups {
		if !layer.Visible {
			continue
		}
		if err := r.RenderObjectGroup(i); err != nil {
			return err
		}
	}
	return nil
}

// RenderObjectGroup renders a single object group
func (r *Renderer) RenderObjectGroup(i int) error {
	if i >= len(r.m.ObjectGroups) {
		return ErrOutOfBounds
	}

	layer := r.m.ObjectGroups[i]
	return r._renderObjectGroup(layer)
}

func (r *Renderer) _renderObjectGroup(objectGroup *tiled.ObjectGroup) error {
	objs := objectGroup.Objects

	// sort objects from left top to right down
	objs = utils.SortAnySlice(objs, func(a, b *tiled.Object) bool {
		if a.Y != b.Y {
			return a.Y < b.Y
		}

		return a.X < b.X
	})

	for _, obj := range objs {
		if err := r.renderOneObject(objectGroup, obj); err != nil {
			return err
		}
	}
	return nil
}

// RenderGroupObjectGroup renders single object group in a certain group.
func (r *Renderer) RenderGroupObjectGroup(groupIdx, objectGroupId int) error {
	if groupIdx >= len(r.m.Groups) {
		return ErrOutOfBounds
	}

	group := r.m.Groups[groupIdx]

	if objectGroupId >= len(group.ObjectGroups) {
		return ErrOutOfBounds
	}

	layer := group.ObjectGroups[objectGroupId]
	return r._renderObjectGroup(layer)
}

func (r *Renderer) renderOneObject(layer *tiled.ObjectGroup, o *tiled.Object) error {
	if !o.Visible {
		return nil
	}

	if o.GID == 0 {
		// TODO: o.GID == 0
		return nil
	}

	tile, err := r.m.TileGIDToTile(o.GID)
	if err != nil {
		return err
	}

	img, err := r.getTileImage(tile)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	srcSize := bounds.Size()
	dstSize := image.Pt(int(o.Width), int(o.Height))

	if !srcSize.Eq(dstSize) {
		img = imaging.Resize(img, dstSize.X, dstSize.Y, imaging.NearestNeighbor)
	}

	var originPoint image.Point

	img, originPoint = r._rotateObjectImage(img, o.Rotation)

	bounds = img.Bounds()
	pos := bounds.Add(image.Pt(int(o.X), int(o.Y)).Sub(originPoint))

	if layer.Opacity < 1 {
		mask := image.NewUniform(color.Alpha{uint8(layer.Opacity * 255)})

		draw.DrawMask(r.Result, pos, img, img.Bounds().Min, mask, mask.Bounds().Min, draw.Over)
	} else {
		draw.Draw(r.Result, pos, img, img.Bounds().Min, draw.Over)
	}

	return nil
}

func (r *Renderer) _rotateObjectImage(img image.Image, rotation float64) (newImage image.Image, originPoint image.Point) {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	points := []image.Point{
		image.Pt(0, 0),
		image.Pt(w-1, 0),
		image.Pt(w-1, h-1),
		image.Pt(0, h-1),
	}

	sin, cos := math.Sincos(math.Pi * rotation / 180)

	rotatedPointsX := []float64{}
	rotatedPointsY := []float64{}

	for _, p := range points {
		x := float64(p.X)
		y := float64(p.Y)

		rotatedPointsX = append(rotatedPointsX, x*cos-y*sin)
		rotatedPointsY = append(rotatedPointsY, x*sin+y*cos)
	}

	rotatedMinX := rotatedPointsX[0]
	rotatedMinY := rotatedPointsY[0]

	for i := 1; i < 4; i++ {
		rotatedMinX = math.Min(rotatedMinX, rotatedPointsX[i])
		rotatedMinY = math.Min(rotatedMinY, rotatedPointsY[i])
	}

	originPoint = image.Pt(int(rotatedPointsX[3]-rotatedMinX), int(rotatedPointsY[3]-rotatedMinY))

	return imaging.Rotate(img, -rotation, color.RGBA{}), originPoint
}
