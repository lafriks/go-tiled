/*
Copyright (c) 2026 Lauris Bukšis <lauris@nix.lv>

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

// StaggeredRendererEngine represents a staggered rendering engine, supporting
// both StaggerAxis "x" and "y" with either StaggerIndex.
//
// Despite the name, orientation="staggered" is not a hexagonal map -- it's
// what Tiled's own "New Map" dialog calls "Isometric (Staggered)": ordinary
// (non-hex) diamond-ish tile art arranged in a staggered/offset grid to fake
// an isometric look without true isometric projection. "Hexagonal (Staggered)"
// is the separate orientation="hexagonal", for actual hex tiles.
//
// The two share an implementation in Tiled itself, though: hex support was
// added later by generalizing the pre-existing isometric-staggered renderer
// with a HexSideLength (see
// https://discourse.mapeditor.org/t/support-for-hexagonal-maps-added/47), and
// Tiled's StaggeredRenderer still literally inherits from HexagonalRenderer
// today with no changes to tile placement or bounding-box math (see
// https://github.com/mapeditor/tiled/blob/master/src/libtiled/staggeredrenderer.h).
// The "staggered" orientation has no hexsidelength attribute at all, so
// tiled.Map.HexSideLength is naturally its zero value for these maps, which
// reduces HexagonalRendererEngine's geometry to exactly this staggered,
// half-tile brick-like offset grid -- so this type only needs to embed it.
type StaggeredRendererEngine struct {
	HexagonalRendererEngine
}
