package doublearray

import (
	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/pkg/errors"
)

type builder struct {
	factory nodeFactory
}

func newBuilder() *builder {
	return &builder{
		factory: &fatFactory{},
	}
}

func (b *builder) FromBytes(xs [][]byte) (*DoubleArray, error) {
	ret := &DoubleArray{
		nodes: make([]node.Node, len(xs)*2),
	}
	ks := keyset.FromBytes(xs)
	if err := b.build(ret, ks); err != nil {
		return nil, err
	}
	return ret, nil
}

func (b *builder) build(da *DoubleArray, ks keyset.KeySet) error {
	b.init(da, 0)
	ks.Sort()
	return ks.Walk(func(prefix word.Word, branch []word.Code, vals []uint32) error {
		return b.insert(da, prefix, branch, vals)
	})
}

func (b *builder) insert(da *DoubleArray, prefix word.Word, branch []word.Code, vals []uint32) error {
	//log.Printf("insert(prefix, branch)=(%v, %v)", prefix, branch)

	if err := b.checkBranch(branch); err != nil {
		return err
	}

	// prefixが入っているところを探して、
	index, err := da.getIndex(prefix)
	if err != nil {
		return err
	}

	// branch を全部格納できるoffsetを探して、
	offset := b.findValidOffset(da, branch)

	maxIndex := offset.Forward(branch[len(branch)-1])
	b.ensure(da, maxIndex)

	// nodes[at]にbranchを格納できるoffsetを指定
	da.nodes[index].SetOffset(offset)

	prev := word.NONE
	for i, c := range branch {
		// branch には同じラベルの枝が複数あることがある
		// もう追加していたらスキップ
		if c == prev {
			continue
		}
		next := offset.Forward(c)
		b.popNode(da, next)
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

func (b *builder) checkBranch(branch []word.Code) error {
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

func (b *builder) init(da *DoubleArray, after int) {
	if after == 0 {
		da.nodes[0] = b.factory.root()
		after = 1
	}

	for i := after; i < len(da.nodes); i++ {
		da.nodes[i] = b.factory.node(i)
	}
}

func (b *builder) extend(da *DoubleArray) {
	max := len(da.nodes)
	da.nodes = append(da.nodes, make([]node.Node, len(da.nodes))...)
	b.init(da, max)
}

func (b *builder) ensure(da *DoubleArray, i node.Index) {
	for len(da.nodes) <= int(i) {
		b.extend(da)
	}
}

func (b *builder) popNode(da *DoubleArray, i node.Index) {
	// これから nodes[i] を使うための準備
	// nodes[i] を prev/next にしているnodeから node[i]を取り除く

	prev := da.nodes[i].GetPrevEmptyNode()
	next := da.nodes[i].GetNextEmptyNode()

	// next にアクセスできるように、必要があれば拡張
	b.ensure(da, next)

	// 1. nodes[i].prev の next に nodes[i].next を繋ぐ
	da.nodes[prev].SetNextEmptyNode(next)
	// 2. nodes[i].next の prev に nodes[i].prev を繋ぐ
	da.nodes[next].SetPrevEmptyNode(prev)
}

func (b *builder) findValidOffset(da *DoubleArray, cs word.Word) node.Index {
	index, offset := b.findOffset(da, da.nodes[0].GetNextEmptyNode(), cs[0])

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
			// 使用済みだった
			// 次の未使用ノードを試す
			index, offset = b.findOffset(da, da.nodes[index].GetNextEmptyNode(), cs[0])
			// cs[0] からやりなおし
			i = 0
		}
	}
	return offset
}
func (b *builder) findOffset(da *DoubleArray, index node.Index, branch word.Code) (node.Index, node.Index) {
	for {
		offset, err := index.Backward(branch)
		if err == nil {
			return index, offset
		}
		b.ensure(da, index)
		index = da.nodes[index].GetNextEmptyNode()
	}
}
