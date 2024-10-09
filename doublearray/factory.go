package doublearray

import (
	"github.com/ajiyoshi-vg/external/emit"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/result"
)

type Factory struct {
	ch   chan item.Item
	done <-chan result.Result[*DoubleArray]
}

func NewFactory(b *Builder) *Factory {
	ch := make(chan item.Item)
	done := factory(b, ch)

	ret := &Factory{
		ch:   ch,
		done: done,
	}

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

func factory(b *Builder, ch <-chan item.Item) <-chan result.Result[*DoubleArray] {
	done := make(chan result.Result[*DoubleArray])
	go func() {
		defer close(done)
		x, err := b.StreamBuild(emit.Chan(ch))
		if err != nil {
			done <- result.NG[*DoubleArray](err)
		} else {
			done <- result.OK(x)
		}
	}()
	return done
}
