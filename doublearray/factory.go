package doublearray

import (
	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/result"
)

type Factory struct {
	ch   chan item.Item
	done <-chan result.Result[*DoubleArray]
}

func NewFactory(b *Builder) *Factory {
	ch := make(chan item.Item)

	return &Factory{
		ch:   ch,
		done: factory(b, ch),
	}
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
		done <- result.New(b.StreamBuild(scan.Chan(ch)))
	}()
	return done
}
