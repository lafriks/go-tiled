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

package utils

import (
	"sort"

	"github.com/lafriks/go-tiled"
)

func SortObjectsLess(a, b *tiled.Object) bool {
	if a.Y != b.Y {
		return a.Y < b.Y
	}

	return a.X < b.X
}

func SortAny[T any](data []T, lessMethod func(a, b T) bool) []T {
	s := &Sortable[T]{
		Data:       data,
		LessMethod: lessMethod,
	}
	sort.Sort(s)
	return s.Data
}

type Sortable[T any] struct {
	Data       []T
	LessMethod func(a, b T) bool
}

func (s *Sortable[T]) Swap(i, j int) {
	tmp := (s.Data)[i]
	(s.Data)[i] = (s.Data)[j]
	(s.Data)[j] = tmp
}

func (s *Sortable[T]) Less(i, j int) bool {
	return s.LessMethod(s.Data[i], s.Data[j])
}

func (s *Sortable[T]) Len() int {
	return len(s.Data)
}
