/*
   Copyright 2021 Google LLC.
   SPDX-License-Identifier: Apache-2.0
*/

package tiled_test

import (
	"testing"

	"github.com/lafriks/go-tiled"
)

func TestLoadTilesetReferencedOnlyByObjectGroup(t *testing.T) {
	const mapPath = "assets/test_tileobject.tmx"

	m, err := tiled.LoadFile(mapPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, ts := range m.Tilesets {
		if !ts.SourceLoaded {
			t.Errorf("Tileset %q was not loaded by map %q", ts.Source, mapPath)
		}
	}
}
