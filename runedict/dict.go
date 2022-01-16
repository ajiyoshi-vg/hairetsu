package runedict

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

type RuneDict map[rune]word.Code

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

func (d RuneDict) MarshalText() (string, error) {
	rs := make([]record, 0, len(d))
	for r, c := range d {
		r := record{Code: c, Rune: r, Char: fmt.Sprintf("%c", r)}
		rs = append(rs, r)
	}
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Code < rs[j].Code
	})
	ss := make([]string, 0, len(d))
	for _, r := range rs {
		buf, err := json.Marshal(&r)
		if err != nil {
			return "", err
		}
		ss = append(ss, string(buf))
	}
	return strings.Join(ss, "\n"), nil
}

func (d RuneDict) UnmarshalText(s string) error {
	scan := bufio.NewScanner(bytes.NewBufferString(s))
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
