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

// Properties wraps any number of custom properties
type Properties []*Property

// Property is used for custom properties
type Property struct {
	// The name of the property.
	Name string `xml:"name,attr"`
	// The type of the property. Can be string (default), int, float, bool, color or file (since 0.16, with color and file added in 0.17).
	Type string `xml:"type,attr"`
	// The value of the property.
	// Boolean properties have a value of either "true" or "false".
	// Color properties are stored in the format #AARRGGBB.
	// File properties are stored as paths relative from the location of the map file.
	Value string `xml:"value,attr"`
}

// Get finds all properties by specified name
func (p Properties) Get(name string) []string {
	values := make([]string, 0)
	for i := 0; i < len(p); i++ {
		if p[i].Name == name {
			values = append(values, p[i].Value)
		}
	}

	return values
}

// GetString finds first string property by specified name
func (p Properties) GetString(name string) string {
	v := ""
	for i := 0; i < len(p); i++ {
		if p[i].Name == name {
			if p[i].Type == "" {
				return p[i].Value
			} else if v == "" {
				v = p[i].Value
			}
		}
	}

	return v
}

// GetBool finds first bool property by specified name
func (p Properties) GetBool(name string) bool {
	for i := 0; i < len(p); i++ {
		if p[i].Name == name && p[i].Type == "Boolean" {
			return p[i].Value == "true"
		}
	}

	return p.GetString(name) == "true"
}
