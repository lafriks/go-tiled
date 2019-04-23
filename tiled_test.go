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
	"bytes"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetAssetsDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "assets")
}

func TestLoadFromReader(t *testing.T) {
	r := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<map version="1.2" tiledversion="1.2.1" orientation="orthogonal" renderorder="right-down" width="4" height="4" tilewidth="16" tileheight="16" infinite="0" nextlayerid="2" nextobjectid="2">
<layer id="1" name="Tile Layer 1" width="4" height="4">
<data encoding="csv">
0,0,0,0,
0,0,0,0,
0,0,0,0,
0,0,0,0
</data>
</layer>
</map>`)
	m, err := LoadFromReader(GetAssetsDirectory(), r)

	assert.NoError(t, err)
	assert.NotNil(t, m)
	assert.Len(t, m.Layers, 1)
}

func TestLoadFromReaderError(t *testing.T) {
	r := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<map version="1.2" tiledversion="1.2.1" orientation="orthogonal" renderorder="right-down" width="4" height="4" tilewidth="16" tileheight="16" infinite="0" nextlayerid="2" nextobjectid="2">
<layer id="1" name="Tile Layer 1" width="4" height="4">
<data encoding="csv">
</layer>
</map>`)
	m, err := LoadFromReader(GetAssetsDirectory(), r)

	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestLoadFromFile(t *testing.T) {
	m, err := LoadFromFile(filepath.Join(GetAssetsDirectory(), "test.tmx"))

	assert.NoError(t, err)
	assert.NotNil(t, m)
	assert.Len(t, m.Layers, 1)
}

func TestLoadFromFileError(t *testing.T) {
	m, err := LoadFromFile(filepath.Join(GetAssetsDirectory(), "invalid.tmx"))

	assert.Error(t, err)
	assert.Nil(t, m)
}
