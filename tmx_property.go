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
	"encoding/hex"
	"encoding/xml"
	"image/color"
	"strconv"
)

// Properties wraps any number of custom properties
type Properties []*Property

// Property is used for custom properties
type Property struct {
	// The name of the property.
	Name string
	// The type of the property. Can be string (default), int, float, bool, color or file (since 0.16, with color and file added in 0.17).
	Type string
	// The value of the property.
	// Boolean properties have a value of either "true" or "false".
	// Color properties are stored in the format #AARRGGBB.
	// File properties are stored as paths relative from the location of the map file.
	Value string
}

// UnmarshalXML implements the xml.Unmarshaler interface for Property. Setting Value even if it's in the inner text.
func (p *Property) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var valueFoundInAttr bool
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "name":
			p.Name = attr.Value
		case "type":
			p.Type = attr.Value
		case "value":
			p.Value = attr.Value
			valueFoundInAttr = true
		}
	}
	if valueFoundInAttr {
		return d.Skip()
	}

	var innerText string
	if err := d.DecodeElement(&innerText, &start); err != nil {
		return err
	}
	p.Value = innerText

	return nil
}

// Get finds all properties by specified name
func (p Properties) Get(name string) []string {
	var values []string
	for _, property := range p {
		if property.Name == name {
			values = append(values, property.Value)
		}
	}
	return values
}

// GetString finds first string property by specified name
func (p Properties) GetString(name string) string {
	var v string
	for _, property := range p {
		if property.Name == name {
			if property.Type == "" {
				return property.Value
			} else if v == "" {
				v = property.Value
			}
		}
	}
	return v
}

// GetBool finds first bool property by specified name
func (p Properties) GetBool(name string) bool {
	for _, property := range p {
		if property.Name == name && property.Type == "boolean" {
			return property.Value == "true"
		}
	}
	return p.GetString(name) == "true"
}

// GetInt finds first int property by specified name
func (p Properties) GetInt(name string) int {
	for _, property := range p {
		if property.Name == name && property.Type == "int" {
			v, err := strconv.Atoi(property.Value)
			if err != nil {
				continue
			}
			return v
		}
	}
	return 0
}

// GetFloat finds first float property by specified name
func (p Properties) GetFloat(name string) float64 {
	for _, property := range p {
		if property.Name == name && property.Type == "float" {
			v, err := strconv.ParseFloat(property.Value, 64)
			if err != nil {
				continue
			}
			return v
		}
	}
	return 0
}

// GetColor returns a color.Color by parsing the first property found using
// name. If unable to parse the value or find the value nil is returned.
func (p Properties) GetColor(name string) color.Color {
	for _, property := range p {
		if property.Name != name || property.Type != "color" {
			continue
		}
		str := property.Value
		if len(str) < 6 {
			continue
		}
		if str[0] == '#' {
			str = str[1:]
		}
		b, err := hex.DecodeString(str)
		if err != nil {
			continue
		}
		return &color.RGBA{b[1], b[2], b[3], b[0]}
	}
	return nil
}
