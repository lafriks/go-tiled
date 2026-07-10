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
	// Raw data
	RawData []byte `xml:",innerxml"`
	// Only used when layer encoding is xml
	DataTiles []DataTile `xml:"tile"`
}

// DataTile defines the value of a single tile on a tile layer
type DataTile struct {
	// The global tile ID.
	GID uint32 `xml:"gid,attr"`
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
