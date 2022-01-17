package runedict

import (
	"bufio"
	"bytes"
	"encoding"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"sort"

	"github.com/ajiyoshi-vg/hairetsu/lines"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type RuneDict map[rune]word.Code

var (
	_ encoding.BinaryMarshaler   = RuneDict(nil)
	_ encoding.BinaryUnmarshaler = RuneDict(nil)
	_ encoding.TextMarshaler     = RuneDict(nil)
	_ encoding.TextUnmarshaler   = RuneDict(nil)
)

type Builder struct {
	runeCount map[rune]uint32
}

func New(ss []string) RuneDict {
	b := NewBuilder()
	for _, s := range ss {
		b.Add(s)
	}
	return b.Build()
}

func (d RuneDict) Code(r rune) word.Code {
	ret, ok := d[r]
	if !ok {
		return word.NONE
	}
	return ret
}

func (d RuneDict) Word(s string) (word.Word, error) {
	ret := make(word.Word, 0, len(s))
	for _, r := range s {
		c, ok := d[r]
		if !ok {
			return nil, fmt.Errorf("unknown rune(%c)", r)
		} else {
			ret = append(ret, c)
		}
	}
	return ret, nil
}

type record struct {
	Code word.Code `json:"code"`
	Rune rune      `json:"rune"`
	Char string    `json:"char"`
}

func (d RuneDict) MarshalText() ([]byte, error) {
	rs := make([]record, 0, len(d))
	for r, c := range d {
		r := record{Code: c, Rune: r, Char: fmt.Sprintf("%c", r)}
		rs = append(rs, r)
	}
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Code < rs[j].Code
	})
	buf := &bytes.Buffer{}
	for _, r := range rs {
		x, err := json.Marshal(&r)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(x)
		if err != nil {
			return nil, err
		}
		_, err = buf.WriteRune('\n')
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (d RuneDict) UnmarshalText(s []byte) error {
	scan := bufio.NewScanner(bytes.NewBuffer(s))
	tmp := &record{}

	for i := 0; scan.Scan(); i++ {
		line := scan.Text()
		if err := json.Unmarshal([]byte(line), &tmp); err != nil {
			return err
		}
		d[tmp.Rune] = tmp.Code
	}

	return nil
}

func (d RuneDict) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	for r, c := range d {
		if err := binary.Write(buf, binary.BigEndian, uint32(r)); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, uint32(c)); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (d RuneDict) UnmarshalBinary(s []byte) error {
	buf := bytes.NewReader(s)
	for {
		var r rune
		var c word.Code
		err := binary.Read(buf, binary.BigEndian, &r)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		if err := binary.Read(buf, binary.BigEndian, &c); err != nil {
			return err
		}
		d[r] = c
	}
}

func (d RuneDict) WriteTo(w io.Writer) (int64, error) {
	out := bufio.NewWriter(w)
	defer out.Flush()

	buf, err := d.MarshalBinary()
	if err != nil {
		return 0, err
	}
	n, err := out.Write(buf)
	return int64(n), err
}

func (d RuneDict) ReadFrom(r io.Reader) (int64, error) {
	buf, err := ioutil.ReadAll(r)
	ret := int64(len(buf))
	if err != nil {
		return ret, err
	}
	if err := d.UnmarshalBinary(buf); err != nil {
		return ret, err
	}
	return ret, nil
}

func NewBuilder() *Builder {
	return &Builder{
		runeCount: map[rune]uint32{},
	}
}

func (b *Builder) Add(s string) {
	for _, r := range s {
		b.runeCount[r] += 1
	}
}

func (b *Builder) Build() RuneDict {
	type tmp struct {
		r rune
		n uint32
	}

	buf := make([]tmp, 0, len(b.runeCount))
	for r, n := range b.runeCount {
		buf = append(buf, tmp{r: r, n: n})
	}

	sort.Slice(buf, func(i, j int) bool {
		return buf[i].n > buf[j].n
	})

	ret := make(RuneDict, len(buf))
	for i, x := range buf {
		ret[x.r] = word.Code(i)
	}
	return ret
}

func FromLines(r io.Reader) (RuneDict, error) {
	b := NewBuilder()
	lines.AsString(r, func(s string) error {
		b.Add(s)
		return nil
	})
	return b.Build(), nil
}
