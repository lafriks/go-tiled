/*
Copyright (c) 2017 Lauris Bukšis <lauris@nix.lv>

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
func (r *Renderer) RenderGroup(groupID int) error {
	if groupID >= len(r.m.Groups) {
		return ErrOutOfBounds
	}

	group := r.m.Groups[groupID]
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

// positionedObject pairs an object with its already-projected screen position,
// so sorting and rendering both use one, orientation-correct, source of truth.
type positionedObject struct {
	obj    *tiled.Object
	sx, sy float64
}

func (r *Renderer) _renderObjectGroup(objectGroup *tiled.ObjectGroup) error {
	objs := make([]positionedObject, len(objectGroup.Objects))
	for i, obj := range objectGroup.Objects {
		sx, sy := r.engine.PixelToScreenCoords(obj.X, obj.Y)
		objs[i] = positionedObject{obj: obj, sx: sx, sy: sy}
	}

	// sort objects from screen top to screen bottom, so they draw back-to-front;
	// raw object X/Y aren't screen order for orientations like isometric, where the
	// grid is sheared, so this must sort by the projected position, not o.X/o.Y.
	objs = utils.SortAnySlice(objs, func(a, b positionedObject) bool {
		if a.sy != b.sy {
			return a.sy < b.sy
		}

		return a.sx < b.sx
	})

	for _, obj := range objs {
		if err := r.renderOneObject(objectGroup, obj.obj, obj.sx, obj.sy); err != nil {
			return err
		}
	}
	return nil
}

// RenderGroupObjectGroup renders single object group in a certain group.
func (r *Renderer) RenderGroupObjectGroup(groupID, objectGroupID int) error {
	if groupID >= len(r.m.Groups) {
		return ErrOutOfBounds
	}

	group := r.m.Groups[groupID]

	if objectGroupID >= len(group.ObjectGroups) {
		return ErrOutOfBounds
	}

	layer := group.ObjectGroups[objectGroupID]
	return r._renderObjectGroup(layer)
}

func (r *Renderer) renderOneObject(layer *tiled.ObjectGroup, o *tiled.Object, screenX, screenY float64) error {
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

	anchor := r.engine.GetObjectAnchor(img.Bounds().Size())

	var originPoint image.Point

	img, originPoint = r._rotateObjectImage(img, o.Rotation, anchor)

	bounds = img.Bounds()
	pos := bounds.Add(image.Pt(int(screenX), int(screenY)).Sub(originPoint))

	if layer.Opacity < 1 {
		mask := image.NewUniform(color.Alpha{uint8(layer.Opacity * 255)})

		draw.DrawMask(r.Result, pos, img, img.Bounds().Min, mask, mask.Bounds().Min, draw.Over)
	} else {
		draw.Draw(r.Result, pos, img, img.Bounds().Min, draw.Over)
	}

	return nil
}

// _rotateObjectImage rotates img around anchor (given in img's own, unrotated
// coordinate space) and returns the rotated image along with anchor's new
// position within it, so the caller can re-align the rotated image to the
// same screen point that anchor represented before rotation.
func (r *Renderer) _rotateObjectImage(img image.Image, rotation float64, anchor image.Point) (newImage image.Image, originPoint image.Point) {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	corners := []image.Point{
		image.Pt(0, 0),
		image.Pt(w-1, 0),
		image.Pt(w-1, h-1),
		image.Pt(0, h-1),
	}

	sin, cos := math.Sincos(math.Pi * rotation / 180)

	rotate := func(p image.Point) (float64, float64) {
		x := float64(p.X)
		y := float64(p.Y)
		return x*cos - y*sin, x*sin + y*cos
	}

	rx0, ry0 := rotate(corners[0])
	rotatedMinX, rotatedMinY := rx0, ry0

	for _, c := range corners[1:] {
		rx, ry := rotate(c)
		rotatedMinX = math.Min(rotatedMinX, rx)
		rotatedMinY = math.Min(rotatedMinY, ry)
	}

	anchorX, anchorY := rotate(anchor)
	originPoint = image.Pt(int(anchorX-rotatedMinX), int(anchorY-rotatedMinY))

	return imaging.Rotate(img, -rotation, color.RGBA{}), originPoint
}
