package stream

import (
	"iter"
	"log"

	"github.com/ajiyoshi-vg/external"
	"github.com/ajiyoshi-vg/external/emit"
)

func Sort[T any](seq iter.Seq[T], cmp func(T, T) int, opt ...external.Option) (iter.Seq[T], int, error) {
	s := external.NewSplitter(cmp, opt...)
	chunk, err := s.Split(seq)
	if err != nil {
		return nil, 0, err
	}
	ss, err := chunk.Iters()
	if err != nil {
		return nil, 0, err
	}
	m := external.NewMerger(cmp)
	sorted := m.Merge(ss)
	return func(yield func(T) bool) {
		defer func() {
			if err := chunk.Clean(); err != nil {
				log.Println(err)
			}
		}()
		emit.All(sorted, yield)
	}, chunk.Length(), nil
}
