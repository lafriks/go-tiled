package tiled

import (
	"encoding/hex"
	"encoding/xml"
	"errors"
	"image/color"
	"strings"
)

// HexColor handles the conversion between hex color strings in form #AARRGGBB
// to color.RGBA structure. Be aware that this doesn't match CSS hex color because
// the alpha channel appears first.
type HexColor struct {
	c color.RGBA
}

// ParseHexColor converts hex color strings into HexColor structures
// This function can handle colors with and withouth optional alpha channel
// The leading '#' character is not required for backwards compatibility with Transparency Tiled filed
// https://doc.mapeditor.org/en/stable/reference/tmx-map-format/#image
func ParseHexColor(s string) (HexColor, error) {
	c, err := parseHexColor(s)
	if err != nil {
		return HexColor{}, err
	}
	return HexColor{c: c}, nil
}

// NewHexColor is a shorthand to build a HexColor
func NewHexColor(r, g, b, a uint32) HexColor {
	return HexColor{
		c: color.RGBA{
			uint8(r),
			uint8(g),
			uint8(b),
			uint8(a),
		},
	}
}

// RGBA implements color.Color interface
func (color *HexColor) RGBA() (r, g, b, a uint32) {
	return color.c.RGBA()
}

func (color *HexColor) String() string {
	src := []byte{
		color.c.A,
		color.c.R,
		color.c.G,
		color.c.B,
	}
	if color.c.A == 255 {
		src = src[1:]
	}

	dst := make([]byte, hex.EncodedLen(len(src))+1)
	hex.Encode(dst[1:], src)
	dst[0] = '#'
	return string(dst)
}

// UnmarshalXMLAttr implements xml.UnmarshalerAttr
func (color *HexColor) UnmarshalXMLAttr(attr xml.Attr) error {
	c, err := parseHexColor(attr.Value)
	if err != nil {
		return err
	}
	color.c = c
	return nil
}

// MarshalXMLAttr implements xml.MarshalerAttr
func (color *HexColor) MarshalXMLAttr(name xml.Name) (attr xml.Attr, err error) {
	attr.Name = name
	if color != nil {
		attr.Value = color.String()
	}
	return
}

func parseHexColor(s string) (c color.RGBA, err error) {
	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errors.New("Invalid Format")
		return 0
	}

	s = strings.TrimPrefix(s, "#")

	switch len(s) {
	case 8:
		c.A = hexToByte(s[0])<<4 + hexToByte(s[1])
		c.R = hexToByte(s[2])<<4 + hexToByte(s[3])
		c.G = hexToByte(s[4])<<4 + hexToByte(s[5])
		c.B = hexToByte(s[6])<<4 + hexToByte(s[7])
	case 6:
		c.A = 0xff
		c.R = hexToByte(s[0])<<4 + hexToByte(s[1])
		c.G = hexToByte(s[2])<<4 + hexToByte(s[3])
		c.B = hexToByte(s[4])<<4 + hexToByte(s[5])
	case 4:
		c.A = hexToByte(s[0]) * 17
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	case 3:
		c.A = 0xff
		c.R = hexToByte(s[0]) * 17
		c.G = hexToByte(s[1]) * 17
		c.B = hexToByte(s[2]) * 17
	default:
		err = errors.New("Invalid Format")
	}
	return
}
