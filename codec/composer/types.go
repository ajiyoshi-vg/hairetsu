package composer

import (
	"github.com/ajiyoshi-vg/hairetsu/codec/bytes"
	"github.com/ajiyoshi-vg/hairetsu/codec/runes"
	"github.com/ajiyoshi-vg/hairetsu/codec/u16s"
)

func NewInt16[D u16s.WordDict](d D) *Composer[
	D,
	*u16s.Encoder[D],
	[]byte,
	uint16,
] {
	return NewComposer(
		d,
		u16s.NewEncoder[D],
		u16s.FromReadSeeker,
	)
}
func NewBytes[D bytes.WordDict](d D) *Composer[
	D,
	*bytes.Encoder[D],
	[]byte,
	byte,
] {
	return NewComposer(
		d,
		bytes.NewEncoder[D],
		bytes.FromReadSeeker,
	)
}

func NewRunes[D runes.WordDict](d D) *Composer[
	D,
	*runes.Encoder[D],
	string,
	rune,
] {
	return NewComposer(
		d,
		runes.NewEncoder[D],
		runes.FromReadSeeker,
	)
}

var (
	_ = NewInt16(u16s.NewMapDict())
	_ = NewInt16(u16s.NewArrayDict())
	_ = NewInt16(u16s.NewIdentityDict())
	_ = NewBytes(bytes.NewMapDict())
	_ = NewBytes(bytes.NewArrayDict())
	_ = NewBytes(bytes.NewIdentityDict())
	_ = NewRunes(runes.NewMapDict())
	_ = NewRunes(runes.NewIdentityDict())
)
