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
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	// ErrUnknownCompression error is returned when file contains invalid compression method
	ErrUnknownCompression = errors.New("tiled: invalid compression method")
)

// Data is raw data
type Data struct {
	// The encoding used to encode the tile layer data. When used, it can be "base64" and "csv" at the moment.
	Encoding string `xml:"encoding,attr"`
	// The compression used to compress the tile layer data. Tiled Qt supports "gzip" and "zlib".
	Compression string `xml:"compression,attr"`
	// Raw data
	RawData []byte `xml:",innerxml"`
	// Only used when layer encoding is xml
	DataTiles []*DataTile `xml:"tile"`
}

// DataTile defines the value of a single tile on a tile layer
type DataTile struct {
	// The global tile ID.
	GID uint32 `xml:"gid,attr"`
}

func (d *Data) decodeBase64() (data []byte, err error) {
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

	return ioutil.ReadAll(comr)
}

func (d *Data) decodeCSV() ([]uint32, error) {
	cleaner := func(r rune) rune {
		if (r >= '0' && r <= '9') || r == ',' {
			return r
		}
		return -1
	}
	rawDataClean := strings.Map(cleaner, string(d.RawData))

	str := strings.Split(string(rawDataClean), ",")

	gids := make([]uint32, len(str))
	for i, s := range str {
		var id uint64
		var err error
		if id, err = strconv.ParseUint(s, 10, 32); err != nil {
			return nil, err
		}
		gids[i] = uint32(id)
	}
	return gids, nil
}
