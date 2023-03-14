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
