package doublearray

import (
	"fmt"
	"io"
	"log"

	"github.com/ajiyoshi-vg/hairetsu/keytree"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Builder struct {
	progress Progress
}

type Progress interface {
	SetMax(int)
	Add(int) error
}

func NewBuilder(opt ...Option) *Builder {
	ret := &Builder{}
	for _, f := range opt {
		f(ret)
	}
	return ret
}

func (b *Builder) ReadFrom(da *DoubleArray, r io.Reader) (int64, error) {
	var ret int64
	length := 8
	buf := make([]byte, length)
	for i := node.Index(0); ; i++ {
		n, err := r.Read(buf)
		ret += int64(n)

		if n == length {
			b.ensure(da, i)
			if err := da.nodes[i].UnmarshalBinary(buf); err != nil {
				return ret, err
			}
		}

		if n > 0 && n < len(buf) {
			return ret, fmt.Errorf("short read(%d), bad align at %d", n, ret)
		}

		if io.EOF == err {
			return ret, nil
		}
	}
}

type Walker interface {
	WalkNode(func(word.Word, []word.Code, *uint32) error) error
	WalkLeaf(func(word.Word, uint32) error) error
	LeafNum() int
}

var (
	_ Walker = (*keytree.Tree)(nil)
)

func (b *Builder) Build(da *DoubleArray, ks Walker) error {
	b.init(da, 0)
	if b.progress != nil {
		b.progress.SetMax(ks.LeafNum())
	}
	return ks.WalkNode(func(prefix word.Word, branch []word.Code, val *uint32) error {
		if val != nil {
			branch = append(branch, word.EOS)
		}
		return b.insert(da, prefix, branch, val)
	})
}

func (b *Builder) insert(da *DoubleArray, prefix word.Word, branch []word.Code, val *uint32) error {
	//logInsert(prefix, branch, val)

	// prefixが入っているところを探して、
	index, err := b.searchIndex(da, prefix)
	if err != nil {
		return err
	}

	// branch を全部格納できるoffsetを探して、
	offset, err := b.findValidOffset(da, branch)
	if err != nil {
		return err
	}

	// nodes[index]にbranchを格納できるoffsetを指定
	da.nodes[index].SetOffset(offset)

	for _, c := range branch {
		next := offset.Forward(c)
		b.ensure(da, next)
		if da.nodes[next].IsUsed() {
			// branch can have same labels
			// skip if it have already inserted
			continue
		}
		if err := b.popNode(da, next); err != nil {
			return err
		}
		da.nodes[next].SetParent(index)

		if c == word.EOS {
			//terminated
			da.nodes[index].Terminate()
			da.nodes[next].SetOffset(node.Index(*val))
			if b.progress != nil {
				b.progress.Add(1)
			}
		}
	}

	return nil
}

func (*Builder) searchIndex(da *DoubleArray, cs word.Word) (node.Index, error) {
	var index node.Index
	for _, c := range cs {
		next := da.nodes[index].GetOffset().Forward(c)
		if int(next) >= len(da.nodes) || !da.nodes[next].IsChildOf(index) {
			return 0, fmt.Errorf("searchIndex(%v) fail", cs)
		}
		index = next
	}
	return index, nil
}

func (b *Builder) findValidOffset(da *DoubleArray, cs word.Word) (node.Index, error) {
	root := node.Index(0)
	index, offset, err := b.findOffset(da, root, cs[0])
	if err != nil {
		return 0, err
	}

	// ensure every offset.Forward(cs[i]) is empty
	for i := 0; i < len(cs); i++ {
		next := offset.Forward(cs[i])

		b.ensure(da, next)

		if da.nodes[next].IsUsed() || next == root {
			// it was used
			index, offset, err = b.findOffset(da, index, cs[0])
			if err != nil {
				return 0, err
			}
			// retry from cs[0]
			i = 0
		}
	}
	return offset, nil
}
func (b *Builder) findOffset(da *DoubleArray, index node.Index, branch word.Code) (node.Index, node.Index, error) {
	next, err := b.nextEmptyNode(da, index)
	if err != nil {
		return 0, 0, err
	}
	offset := next.Backward(branch)
	return next, offset, nil
}

func (b *Builder) init(da *DoubleArray, after int) {
	for i := after; i < len(da.nodes); i++ {
		da.nodes[i].Reset(i)
	}
}
func (b *Builder) extend(da *DoubleArray) {
	max := len(da.nodes)
	da.nodes = append(da.nodes, make([]node.Node, len(da.nodes))...)
	b.init(da, max)
}
func (b *Builder) ensure(da *DoubleArray, i node.Index) {
	for len(da.nodes) <= int(i) {
		b.extend(da)
	}
}
func (b *Builder) popNode(da *DoubleArray, i node.Index) error {
	// prepare to use nodes[i]
	// ensure that nobody keeps nodes[i] as it's prev/next.

	prev, err := b.prevEmptyNode(da, i)
	if err != nil {
		return err
	}
	next, err := b.nextEmptyNode(da, i)
	if err != nil {
		return err
	}

	// 1. let nodes[i].prev.next = nodes[i].next
	if err := b.setNextEmptyNode(da, prev, next); err != nil {
		return err
	}
	// 2. let nodes[i].next.prev = nodes[i].prev
	if err := b.setPrevEmptyNode(da, next, prev); err != nil {
		return err
	}

	return nil
}

func (b *Builder) nextEmptyNode(da *DoubleArray, i node.Index) (node.Index, error) {
	b.ensure(da, i)
	return da.nodes[i].GetNextEmptyNode()
}
func (b *Builder) prevEmptyNode(da *DoubleArray, i node.Index) (node.Index, error) {
	b.ensure(da, i)
	return da.nodes[i].GetPrevEmptyNode()
}
func (b *Builder) setNextEmptyNode(da *DoubleArray, i, next node.Index) error {
	b.ensure(da, i)
	return da.nodes[i].SetNextEmptyNode(next)
}
func (b *Builder) setPrevEmptyNode(da *DoubleArray, i, prev node.Index) error {
	b.ensure(da, i)
	return da.nodes[i].SetPrevEmptyNode(prev)
}

func logInsert(prefix word.Word, branch []word.Code, val *uint32) {
	str := func(x *uint32) string {
		if x == nil {
			return "nil"
		}
		return fmt.Sprintf("%d", *x)
	}
	log.Printf("insert(prefix, branch)=(%v, %v, %s)", prefix, branch, str(val))
}
