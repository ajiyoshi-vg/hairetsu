package doublearray

import (
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/result"
)

type Factory struct {
	ch   chan item.Item
	done chan result.Result[*DoubleArray]
}

func NewFactory(b *Builder) *Factory {
	ret := &Factory{
		ch:   make(chan item.Item),
		done: make(chan result.Result[*DoubleArray]),
	}

	seq := func(yield func(item.Item) bool) {
		for x := range ret.ch {
			if !yield(x) {
				return
			}
		}
	}

	go func() {
		x, err := b.StreamBuild(seq)
		if err != nil {
			ret.done <- result.NG[*DoubleArray](err)
		} else {
			ret.done <- result.OK(x)
		}
	}()

	return ret
}

func (b *Factory) Put(item item.Item) {
	b.ch <- item
}

func (b *Factory) Done() (*DoubleArray, error) {
	close(b.ch)
	ret := <-b.done
	return ret.Result()
}
