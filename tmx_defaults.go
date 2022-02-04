/*
Copyright (c) 2021 Lauris Bukšis-Haberkorns <lauris@nix.lv>
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

type aliasGroup Group
type aliasImageLayer ImageLayer
type internalLayer Layer
type aliasLayer struct {
	internalLayer
	// Layer data in raw format
	Data *Data `xml:"data"`
}
type aliasMap Map
type aliasObject Object
type aliasObjectGroup ObjectGroup
type aliasText Text

// SetDefaults provides default values for Group.
func (a *aliasGroup) SetDefaults() {
	a.Opacity = 1
	a.Visible = true
}

// SetDefaults provides default values for ImageLayer.
func (a *aliasImageLayer) SetDefaults() {
	a.Opacity = 1
	a.Visible = true
}

// SetDefaults provides default values for Layer.
func (a *aliasLayer) SetDefaults() {
	a.internalLayer.Opacity = 1
	a.internalLayer.Visible = true
}

// SetDefaults provides default values for Map.
func (a *aliasMap) SetDefaults() {
	a.RenderOrder = "right-down"
}

// SetDefaults provides default values for Object.
func (a *aliasObject) SetDefaults() {
	a.Visible = b(true)
}

// SetDefaults provides default values for ObjectGroup.
func (a *aliasObjectGroup) SetDefaults() {
	a.Visible = b(true)
	a.Opacity = 1
}

// SetDefaults provides default values for Text.
func (a *aliasText) SetDefaults() {
	a.FontFamily = "sans-serif"
	a.Size = 16
	a.Kerning = b(true)
	a.HAlign = "left"
	a.VAlign = "top"
	a.Color = &HexColor{}
}
