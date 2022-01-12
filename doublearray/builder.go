package doublearray

import (
	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/pkg/errors"
)

type Builder struct {
	factory nodeFactory
}

func NewBuilder() *Builder {
	return &Builder{
		factory: &factory{},
	}
}

func (b *Builder) Build(da *DoubleArray, ks keyset.KeySet) error {
	b.init(da, 0)
	ks.Sort()
	return ks.Walk(func(prefix word.Word, branch []word.Code, vals []uint32) error {
		return b.insert(da, prefix, branch, vals)
	})
}

func (b *Builder) insert(da *DoubleArray, prefix word.Word, branch []word.Code, vals []uint32) error {
	//log.Printf("insert(prefix, branch)=(%v, %v)", prefix, branch)

	if err := b.checkBranch(branch); err != nil {
		return err
	}

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

	maxIndex := offset.Forward(branch[len(branch)-1])
	b.ensure(da, maxIndex)

	// nodes[index]にbranchを格納できるoffsetを指定
	da.nodes[index].SetOffset(offset)

	prev := word.NONE
	for i, c := range branch {
		// branch can have same labels
		// skip if it have already inserted
		if c == prev {
			continue
		}
		next := offset.Forward(c)
		if err := b.popNode(da, next); err != nil {
			return err
		}
		da.nodes[next].SetParent(index)

		if c == word.EOS {
			//終端マーク
			da.nodes[index].Terminate()
			da.nodes[next].SetOffset(node.Index(vals[i]))
		}

		prev = c
	}

	return nil
}

func (*Builder) checkBranch(branch []word.Code) error {
	for i := 1; i < len(branch); i++ {
		if branch[i-1] > branch[i] {
			//branchはソート済みのはずなので後ろに小さな値があったらおかしい
			return errors.Errorf("data is not sorted(%v)", branch)
		}
		if branch[i] == word.EOS {
			//EOSがあるとしたらbranch[0]だけ。途中にあったらおかしい
			return errors.Errorf("bad EOS(%v)", branch)
		}
	}
	return nil
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

	b.ensure(da, next)

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

func (b *Builder) findValidOffset(da *DoubleArray, cs word.Word) (node.Index, error) {
	index, offset, err := b.findOffset(da, 0, cs[0])
	if err != nil {
		return 0, err
	}

	// offset からcs を全部格納可能なところを探す
	for i := 0; i < len(cs); i++ {
		next := offset.Forward(cs[i])

		if int(next) >= len(da.nodes) {
			break
		}

		if int(index) >= len(da.nodes) {
			break
		}

		if da.nodes[next].HasParent() {
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
	if after == 0 {
		da.nodes[0] = b.factory.root()
		after = 1
	}

	for i := after; i < len(da.nodes); i++ {
		da.nodes[i] = b.factory.node(i)
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
