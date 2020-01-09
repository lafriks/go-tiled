/*
Copyright (c) 2020 Lauris Bukšis-Haberkorns <lauris@nix.lv>

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
	"encoding/xml"
)

// Group is a group layer, used to organize the layers of the map in a hierarchy.
type Group struct {
	// Unique ID of the layer.
	// Each layer that added to a map gets a unique id. Even if a layer is deleted,
	// no layer ever gets the same ID. Can not be changed in Tiled. (since Tiled 1.2)
	ID uint32 `xml:"id,attr"`
	// The name of the group layer.
	Name string `xml:"name,attr"`
	// Rendering offset of the image layer in pixels. Defaults to 0. (since 0.15)
	OffsetX int `xml:"offsetx,attr"`
	// Rendering offset of the image layer in pixels. Defaults to 0. (since 0.15)
	OffsetY int `xml:"offsety,attr"`
	// The opacity of the layer as a value from 0 to 1. Defaults to 1.
	Opacity float32 `xml:"opacity,attr"`
	// Whether the layer is shown (1) or hidden (0). Defaults to 1.
	Visible bool `xml:"visible,attr"`
	// Custom properties
	Properties Properties `xml:"properties>property"`
	// Map layers
	Layers []*Layer `xml:"layer"`
	// Map object groups
	ObjectGroups []*ObjectGroup `xml:"objectgroup"`
	// Image layers
	ImageLayers []*ImageLayer `xml:"imagelayer"`
	// Group layers
	Groups []*Group `xml:"group"`
}

// UnmarshalXML decodes a single XML element beginning with the given start element.
func (g *Group) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias Group

	item := Alias{
		Opacity: 1,
		Visible: true,
	}

	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}

	*g = (Group)(item)

	return nil
}
