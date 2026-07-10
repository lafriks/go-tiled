/*
Copyright (c) 2017 Lauris Bukšis <lauris@nix.lv>

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
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"io"
	"math"
	"strconv"
)

// ErrUnknownCompression error is returned when file contains invalid compression method
var ErrUnknownCompression = errors.New("tiled: invalid compression method")

// Data is raw data
type Data struct {
	// The encoding used to encode the tile layer data. When used, it can be "base64" and "csv" at the moment.
	Encoding string `xml:"encoding,attr"`
	// The compression used to compress the tile layer data. Tiled Qt supports "gzip" and "zlib".
	Compression string `xml:"compression,attr"`
	// Raw data. Only populated when Encoding is "csv" or "base64".
	RawData []byte `xml:",innerxml"`
	// Parsed tile elements. Only populated when Encoding is not set.
	DataTiles []DataTile `xml:"tile"`
}

// DataTile defines the value of a single tile on a tile layer
type DataTile struct {
	// The global tile ID.
	GID uint32 `xml:"gid,attr"`
}

// UnmarshalXML implements a hand-rolled decode of the <data> element.
//
// A tile layer's data is either one large text blob (csv/base64) or,
// for the uncompressed per-tile XML encoding, thousands of <tile gid="N"/>
// child elements. Letting encoding/xml decode DataTiles via reflection
// (as driven by the struct tags above) means paying its per-element,
// per-attribute reflection cost for every tile; walking the token stream
// directly and parsing "gid" by hand avoids that and is significantly
// faster for large maps using this encoding. The struct tags are kept so
// Data/DataTile still unmarshal and marshal correctly (round-trip to XML)
// when used directly with encoding/xml outside of this package.
func (d *Data) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "encoding":
			d.Encoding = attr.Value
		case "compression":
			d.Compression = attr.Value
		}
	}

	if d.Encoding != "" {
		return d.decodeCharData(dec)
	}
	return d.decodeTileElements(dec)
}

// decodeCharData collects the element's text content into RawData.
func (d *Data) decodeCharData(dec *xml.Decoder) error {
	var buf bytes.Buffer
	for {
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.CharData:
			buf.Write(t)
		case xml.EndElement:
			d.RawData = buf.Bytes()
			return nil
		}
	}
}

// decodeTileElements parses the <tile gid="N"/> children directly,
// bypassing reflection-based attribute decoding.
func (d *Data) decodeTileElements(dec *xml.Decoder) error {
	for {
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local != "tile" {
				if err := dec.Skip(); err != nil {
					return err
				}
				continue
			}

			var dt DataTile
			for _, attr := range t.Attr {
				if attr.Name.Local != "gid" {
					continue
				}
				v, err := strconv.ParseUint(attr.Value, 10, 32)
				if err != nil {
					return err
				}
				dt.GID = uint32(v)
			}
			d.DataTiles = append(d.DataTiles, dt)

			if err := dec.Skip(); err != nil {
				return err
			}
		case xml.EndElement:
			return nil
		}
	}
}

// decodeBase64 decodes (and, if applicable, decompresses) the raw data.
func (d *Data) decodeBase64(sizeHint int) (data []byte, err error) {
	rawData := bytes.TrimSpace(d.RawData)
	r := bytes.NewReader(rawData)

	encr := base64.NewDecoder(base64.StdEncoding, r)

	var comr io.Reader
	switch d.Compression {
	case "gzip":
		comr, err = gzip.NewReader(encr)
		if err != nil {
			return
		}
	case "zlib":
		comr, err = zlib.NewReader(encr)
		if err != nil {
			return
		}
	case "":
		comr = encr
	default:
		err = ErrUnknownCompression
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, sizeHint))
	if _, err = buf.ReadFrom(comr); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d *Data) decodeCSV() ([]uint32, error) {
	raw := d.RawData

	// Estimate the tile count from the raw byte length
	gids := make([]uint32, 0, len(raw)/2+1)

	var cur uint64
	inNum := false
	for _, c := range raw {
		if c >= '0' && c <= '9' {
			cur = cur*10 + uint64(c-'0')
			if cur > math.MaxUint32 {
				return nil, strconv.ErrRange
			}
			inNum = true
			continue
		}
		if inNum {
			gids = append(gids, uint32(cur))
			cur = 0
			inNum = false
		}
	}
	if inNum {
		gids = append(gids, uint32(cur))
	}

	return gids, nil
}
