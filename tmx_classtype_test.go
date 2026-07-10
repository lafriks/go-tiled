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

package tiled

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These cover Tiled's 1.9 rename of the "type" XML attribute to "class" for
// Object, TilesetTile and WangSet: a pre-1.9 file only has "type", a 1.9+ file
// only has "class" -- both fields should resolve to the same effective value
// regardless of which attribute the source file actually used.

func TestObject_ClassTypeFallback(t *testing.T) {
	var oldFormat Object
	assert.NoError(t, xml.Unmarshal([]byte(`<object id="1" type="npc"/>`), &oldFormat))
	assert.Equal(t, "npc", oldFormat.Type)
	assert.Equal(t, "npc", oldFormat.Class)

	var newFormat Object
	assert.NoError(t, xml.Unmarshal([]byte(`<object id="1" class="npc"/>`), &newFormat))
	assert.Equal(t, "npc", newFormat.Type)
	assert.Equal(t, "npc", newFormat.Class)

	var both Object
	assert.NoError(t, xml.Unmarshal([]byte(`<object id="1" type="old" class="new"/>`), &both))
	assert.Equal(t, "old", both.Type)
	assert.Equal(t, "new", both.Class)
}

func TestTilesetTile_ClassTypeFallback(t *testing.T) {
	var oldFormat TilesetTile
	assert.NoError(t, xml.Unmarshal([]byte(`<tile id="1" type="door"/>`), &oldFormat))
	assert.Equal(t, "door", oldFormat.Type)
	assert.Equal(t, "door", oldFormat.Class)

	var newFormat TilesetTile
	assert.NoError(t, xml.Unmarshal([]byte(`<tile id="1" class="door"/>`), &newFormat))
	assert.Equal(t, "door", newFormat.Type)
	assert.Equal(t, "door", newFormat.Class)
}

func TestWangSet_ClassTypeFallback(t *testing.T) {
	var oldFormat WangSet
	assert.NoError(t, xml.Unmarshal([]byte(`<wangset name="ws" type="corner" tile="-1"/>`), &oldFormat))
	assert.Equal(t, "corner", oldFormat.Type)
	assert.Equal(t, "corner", oldFormat.Class)

	var newFormat WangSet
	assert.NoError(t, xml.Unmarshal([]byte(`<wangset name="ws" class="corner" tile="-1"/>`), &newFormat))
	assert.Equal(t, "corner", newFormat.Type)
	assert.Equal(t, "corner", newFormat.Class)
}

func TestResolveClassType(t *testing.T) {
	tests := []struct {
		name       string
		class, typ string
		wantClass  string
		wantType   string
	}{
		{"both empty", "", "", "", ""},
		{"only class set", "foo", "", "foo", "foo"},
		{"only type set", "", "foo", "foo", "foo"},
		{"both set, distinct", "new", "old", "new", "old"},
		{"both set, same", "foo", "foo", "foo", "foo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClass, gotType := resolveClassType(tt.class, tt.typ)
			assert.Equal(t, tt.wantClass, gotClass)
			assert.Equal(t, tt.wantType, gotType)
		})
	}
}
