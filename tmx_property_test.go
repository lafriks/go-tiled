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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProperty(t *testing.T) {
	props := Properties{
		{
			Name:  "string-name",
			Type:  "string",
			Value: "string-value",
		},
		{
			Name:  "int-name",
			Type:  "int",
			Value: "123",
		},
		{
			Name:  "float-name",
			Type:  "float",
			Value: "1.23",
		},
		{
			Name:  "bool-name",
			Type:  "boolean",
			Value: "true",
		},
	}

	assert.Equal(t, "string-value", props.GetString("string-name"))
	assert.Equal(t, 123, props.GetInt("int-name"))
	assert.Equal(t, 1.23, props.GetFloat("float-name"))
	assert.Equal(t, true, props.GetBool("bool-name"))
}
