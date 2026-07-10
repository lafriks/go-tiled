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
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// buildMapXML generates a synthetic TMX document with a single tile layer of
// the given size and encoding ("csv", "base64" or "" for per-tile XML) so
// benchmarks can exercise realistic, sizeable map data without depending on
// fixture files on disk.
func buildMapXML(b *testing.B, width, height int, encoding, compression string) string {
	b.Helper()

	n := width * height
	gids := make([]uint32, n)
	for i := range gids {
		gids[i] = uint32(i%999) + 1
	}

	var data string
	switch encoding {
	case "csv":
		var sb strings.Builder
		sb.Grow(n * 5)
		for i, gid := range gids {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(strconv.FormatUint(uint64(gid), 10))
		}
		data = fmt.Sprintf(`<data encoding="csv">%s</data>`, sb.String())
	case "base64":
		raw := make([]byte, n*4)
		for i, gid := range gids {
			raw[i*4] = byte(gid)
			raw[i*4+1] = byte(gid >> 8)
			raw[i*4+2] = byte(gid >> 16)
			raw[i*4+3] = byte(gid >> 24)
		}
		var buf bytes.Buffer
		if compression == "zlib" {
			w := zlib.NewWriter(&buf)
			if _, err := w.Write(raw); err != nil {
				b.Fatal(err)
			}
			if err := w.Close(); err != nil {
				b.Fatal(err)
			}
		} else {
			buf.Write(raw)
		}
		encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
		if compression != "" {
			data = fmt.Sprintf(`<data encoding="base64" compression="%s">%s</data>`, compression, encoded)
		} else {
			data = fmt.Sprintf(`<data encoding="base64">%s</data>`, encoded)
		}
	default:
		var sb strings.Builder
		sb.Grow(n * 20)
		for _, gid := range gids {
			fmt.Fprintf(&sb, `<tile gid="%d"/>`, gid)
		}
		data = fmt.Sprintf(`<data>%s</data>`, sb.String())
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<map version="1.2" tiledversion="1.2.4" orientation="orthogonal" renderorder="right-down" width="%d" height="%d" tilewidth="16" tileheight="16" infinite="0" nextlayerid="2" nextobjectid="1">
 <tileset firstgid="1" name="bench" tilewidth="16" tileheight="16" tilecount="1000" columns="10">
  <image source="bench.png" width="160" height="1600"/>
 </tileset>
 <layer id="1" name="Tile Layer 1" width="%d" height="%d">
  %s
 </layer>
</map>`, width, height, width, height, data)
}

func benchmarkLoadReader(b *testing.B, xmlStr string) {
	b.Helper()
	b.SetBytes(int64(len(xmlStr)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := LoadReader("assets", strings.NewReader(xmlStr)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLoadReader_CSV_200x200(b *testing.B) {
	benchmarkLoadReader(b, buildMapXML(b, 200, 200, "csv", ""))
}

func BenchmarkLoadReader_Base64_200x200(b *testing.B) {
	benchmarkLoadReader(b, buildMapXML(b, 200, 200, "base64", ""))
}

func BenchmarkLoadReader_Base64Zlib_200x200(b *testing.B) {
	benchmarkLoadReader(b, buildMapXML(b, 200, 200, "base64", "zlib"))
}

func BenchmarkLoadReader_XMLTiles_200x200(b *testing.B) {
	benchmarkLoadReader(b, buildMapXML(b, 200, 200, "", ""))
}

// BenchmarkLoadFile_Racing exercises the full load pipeline (map + external
// tilesets) against a real, sizeable (100x100, XML tile encoding) fixture.
func BenchmarkLoadFile_Racing(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := LoadFile("assets/racing.tmx"); err != nil {
			b.Fatal(err)
		}
	}
}

// Micro-benchmarks isolating the tile-data decoding hot paths.

func BenchmarkDecodeCSV(b *testing.B) {
	const w, h = 200, 200
	n := w * h
	var sb strings.Builder
	sb.Grow(n * 5)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(i%999 + 1))
	}
	d := &Data{RawData: []byte(sb.String())}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := d.decodeCSV(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeBase64Zlib(b *testing.B) {
	const w, h = 200, 200
	n := w * h
	raw := make([]byte, n*4)
	for i := 0; i < n; i++ {
		gid := uint32(i%999 + 1)
		raw[i*4] = byte(gid)
		raw[i*4+1] = byte(gid >> 8)
		raw[i*4+2] = byte(gid >> 16)
		raw[i*4+3] = byte(gid >> 24)
	}
	var buf bytes.Buffer
	w2 := zlib.NewWriter(&buf)
	if _, err := w2.Write(raw); err != nil {
		b.Fatal(err)
	}
	if err := w2.Close(); err != nil {
		b.Fatal(err)
	}
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	d := &Data{Compression: "zlib", RawData: []byte(encoded)}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := d.decodeBase64(n * 4); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTileGIDToTile(b *testing.B) {
	m := &Map{
		Tilesets: []*Tileset{
			{FirstGID: 1, SourceLoaded: true, TileCount: 1000},
		},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := m.TileGIDToTile(uint32(i%999) + 1); err != nil {
			b.Fatal(err)
		}
	}
}
