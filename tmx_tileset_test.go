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

package tiled

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTileRect(t *testing.T) {
	type Case struct {
		id   uint32
		rect image.Rectangle
	}

	type test struct {
		name  string
		ts    Tileset
		cases []Case
	}

	tests := []test{
		{
			name: "Columns defined",
			ts: Tileset{
				TileCount:  4,
				Columns:    2,
				TileWidth:  10,
				TileHeight: 10,
			},
			cases: []Case{
				{
					id:   0,
					rect: image.Rect(0, 0, 10, 10),
				},
				{
					id:   1,
					rect: image.Rect(0, 0, 10, 10).Add(image.Pt(10, 0)),
				},
				{
					id:   2,
					rect: image.Rect(0, 0, 10, 10).Add(image.Pt(0, 10)),
				},
				{
					id:   3,
					rect: image.Rect(0, 0, 10, 10).Add(image.Pt(10, 10)),
				},
			},
		},
		{
			name: "With Spacing",
			ts: Tileset{
				TileCount:  4,
				Columns:    2,
				TileWidth:  10,
				TileHeight: 10,
				Spacing:    5,
			},
			cases: []Case{
				{
					id:   0,
					rect: image.Rect(0, 0, 10, 10),
				},
				{
					id:   1,
					rect: image.Rect(0, 0, 10, 10).Add(image.Pt(15, 0)),
				},
				{
					id:   2,
					rect: image.Rect(0, 0, 10, 10).Add(image.Pt(0, 15)),
				},
				{
					id:   3,
					rect: image.Rect(0, 0, 10, 10).Add(image.Pt(15, 15)),
				},
			},
		},
		{
			name: "With Margin",
			ts: Tileset{
				TileCount:  4,
				Columns:    2,
				TileWidth:  10,
				TileHeight: 10,
				Margin:     5,
			},
			cases: []Case{
				{
					id:   0,
					rect: image.Rect(5, 5, 15, 15),
				},
				{
					id:   1,
					rect: image.Rect(5, 5, 15, 15).Add(image.Pt(10, 0)),
				},
				{
					id:   2,
					rect: image.Rect(5, 5, 15, 15).Add(image.Pt(0, 10)),
				},
				{
					id:   3,
					rect: image.Rect(5, 5, 15, 15).Add(image.Pt(10, 10)),
				},
			},
		},
		{
			name: "With Margin And Spacing",
			ts: Tileset{
				TileCount:  4,
				Columns:    2,
				TileWidth:  10,
				TileHeight: 10,
				Margin:     5,
				Spacing:    5,
			},
			cases: []Case{
				{
					id:   0,
					rect: image.Rect(5, 5, 15, 15),
				},
				{
					id:   1,
					rect: image.Rect(5, 5, 15, 15).Add(image.Pt(15, 0)),
				},
				{
					id:   2,
					rect: image.Rect(5, 5, 15, 15).Add(image.Pt(0, 15)),
				},
				{
					id:   3,
					rect: image.Rect(5, 5, 15, 15).Add(image.Pt(15, 15)),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, c := range test.cases {
				rect := test.ts.GetTileRect(c.id)
				assert.Equal(t, c.rect, rect)
			}
		})
	}
}
