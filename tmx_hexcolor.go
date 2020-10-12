package tiled

import (
	"encoding/hex"
	"encoding/xml"
	"errors"
	"image/color"
	"strings"
)

type HexColor struct {
	c color.RGBA
}

func ParseHexColor(s string) (HexColor, error) {
	c, err := parseHexColor(s)
	return HexColor{
		c: c,
	}, err

}

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

func (color *HexColor) RGBA() (r, g, b, a uint32) {
	return color.c.RGBA()
}

func (color *HexColor) String() string {
	src := []byte{
		color.c.R,
		color.c.G,
		color.c.B,
		color.c.A,
	}
	if color.c.A == 255 {
		src = src[:len(src)-1]
	}

	dst := make([]byte, hex.EncodedLen(len(src))+1)
	hex.Encode(dst[1:], src)
	dst[0] = '#'
	return string(dst)
}

func (color *HexColor) UnmarshalXMLAttr(attr xml.Attr) error {
	c, err := parseHexColor(attr.Value)
	if err != nil {
		return err
	}
	color.c = c
	return nil
}

func (color *HexColor) MarshalXMLAttr(name xml.Name) (attr xml.Attr, err error) {
	attr.Name = name
	attr.Value = color.String()
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
		c.R = hexToByte(s[0])<<4 + hexToByte(s[1])
		c.G = hexToByte(s[2])<<4 + hexToByte(s[3])
		c.B = hexToByte(s[4])<<4 + hexToByte(s[5])
		c.A = hexToByte(s[6])<<4 + hexToByte(s[7])
	case 6:
		c.R = hexToByte(s[0])<<4 + hexToByte(s[1])
		c.G = hexToByte(s[2])<<4 + hexToByte(s[3])
		c.B = hexToByte(s[4])<<4 + hexToByte(s[5])
		c.A = 0xff
	case 4:
		c.R = hexToByte(s[0]) * 17
		c.G = hexToByte(s[1]) * 17
		c.B = hexToByte(s[2]) * 17
		c.A = hexToByte(s[3]) * 17
	case 3:
		c.R = hexToByte(s[0]) * 17
		c.G = hexToByte(s[1]) * 17
		c.B = hexToByte(s[2]) * 17
		c.A = 0xff
	default:
		err = errors.New("Invalid Format")
	}
	return
}
