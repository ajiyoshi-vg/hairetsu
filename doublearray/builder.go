package doublearray

import (
	"fmt"
	"log"

	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/keytree"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Builder struct {
	factory  nodeFactory
	progress Progress
}

type Progress interface {
	SetMax(int)
	Add(int) error
}

func NewBuilder(opt ...Option) *Builder {
	ret := &Builder{
		factory: &factory{},
	}
	for _, f := range opt {
		f(ret)
	}
	return ret
}

type Walker interface {
	WalkNode(func(word.Word, []word.Code, *uint32) error) error
	WalkLeaf(func(word.Word, uint32) error) error
	LeafNum() int
}

var (
	_ Walker = (*keyset.KeySet)(nil)
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
	index, err := da.searchIndex(prefix)
	if err != nil {
		return err
	}

	// branch を全部格納できるoffsetを探して、
	offset, err := b.findValidOffset(da, branch)
	if err != nil {
		return err
	}

	// nodes[index]にbranchを格納できるoffsetを指定
	da.at(index).SetOffset(offset)

	for _, c := range branch {
		next := offset.Forward(c)
		b.ensure(da, next)
		if da.at(next).IsUsed() {
			// branch can have same labels
			// skip if it have already inserted
			continue
		}
		if err := b.popNode(da, next); err != nil {
			return err
		}
		da.at(next).SetParent(index)

		if val != nil {
			//terminated
			da.at(index).Terminate()
			da.at(next).SetOffset(node.Index(*val))
			if b.progress != nil {
				b.progress.Add(1)
			}
		}
	}

	return nil
}

func (b *Builder) findValidOffset(da *DoubleArray, cs word.Word) (node.Index, error) {
	root := node.Index(0)
	index, offset, err := b.findOffset(da, root, cs[0])
	if err != nil {
		return 0, err
	}

	// offset からcs を全部格納可能なところを探す
	for i := 0; i < len(cs); i++ {
		next := offset.Forward(cs[i])

		b.ensure(da, next)

		if da.at(next).IsUsed() || next == root {
			index, offset, err = b.findOffset(da, index, cs[0])
			if err != nil {
				return 0, err
			}
			// cs[0] からやりなおし
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
	for {
		offset, err := next.Backward(branch)
		if err == nil {
			return next, offset, nil
		}
		next, err = b.nextEmptyNode(da, next)
		if err != nil {
			return 0, 0, err
		}
	}
}

func (b *Builder) init(da *DoubleArray, after int) {
	for i := after; i < len(da.nodes); i++ {
		da.at(node.Index(i)).Reset(i)
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
	return da.at(i).GetNextEmptyNode()
}
func (b *Builder) prevEmptyNode(da *DoubleArray, i node.Index) (node.Index, error) {
	b.ensure(da, i)
	return da.at(i).GetPrevEmptyNode()
}
func (b *Builder) setNextEmptyNode(da *DoubleArray, i, next node.Index) error {
	b.ensure(da, i)
	return da.at(i).SetNextEmptyNode(next)
}
func (b *Builder) setPrevEmptyNode(da *DoubleArray, i, prev node.Index) error {
	b.ensure(da, i)
	return da.at(i).SetPrevEmptyNode(prev)
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
