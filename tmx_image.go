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

// Image source
type Image struct {
	// Used for embedded images, in combination with a data child element. Valid values are file extensions like png, gif, jpg, bmp, etc.
	Format string `xml:"format,attr"`
	// The reference to the tileset image file
	Source string `xml:"source,attr"`
	// Defines a specific color that is treated as transparent (example value: "#FF00FF" for magenta).
	// Up until Tiled 0.12, this value is written out without a # but this is planned to change.
	Trans string `xml:"trans,attr"`
	// The image width in pixels (optional, used for tile index correction when the image changes)
	Width int `xml:"width,attr"`
	// The image height in pixels (optional)
	Height int `xml:"height,attr"`
	// Embeded immage content
	Data *Data `xml:"data,attr"`
}
