package composer

import (
	"github.com/ajiyoshi-vg/hairetsu/codec/bytes"
	"github.com/ajiyoshi-vg/hairetsu/codec/runes"
	"github.com/ajiyoshi-vg/hairetsu/codec/u16s"
)

func NewInt16[D u16s.WordDict](d D) *Composer[
	[]byte,
	uint16,
	D,
	*u16s.Encoder[D],
] {
	return NewComposer(d, u16s.NewBuilder[D]())
}
func NewBytes[D bytes.WordDict](d D) *Composer[
	[]byte,
	byte,
	D,
	*bytes.Encoder[D],
] {
	return NewComposer(d, bytes.NewBuilder[D]())
}

func NewRunes[D runes.WordDict](d D) *Composer[
	string,
	rune,
	D,
	*runes.Encoder[D],
] {
	return NewComposer(d, runes.NewBuilder[D]())
}

var (
	_ Composable[[]byte] = NewInt16(u16s.NewMapDict())
	_ Composable[[]byte] = NewInt16(u16s.NewArrayDict())
	_ Composable[[]byte] = NewInt16(u16s.NewIdentityDict())
	_ Composable[[]byte] = NewBytes(bytes.NewMapDict())
	_ Composable[[]byte] = NewBytes(bytes.NewArrayDict())
	_ Composable[[]byte] = NewBytes(bytes.NewIdentityDict())
	_ Composable[string] = NewRunes(runes.NewMapDict())
	_ Composable[string] = NewRunes(runes.NewIdentityDict())
)
