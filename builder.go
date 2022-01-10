package hairetsu

import (
	"bytes"
	"sort"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/walker"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/pkg/errors"
)

type builder struct {
	id node.Index
}

func (b *builder) FromBytes(xs [][]byte) (*DoubleArray, error) {
	data := make([]word.Word, 0, len(xs))
	for _, x := range xs {
		data = append(data, word.FromBytes(x))
	}
	ret := &DoubleArray{
		nodes:   make([]node.Node, len(xs)*2),
		factory: &fatFactory{},
	}
	if err := b.start(ret, data); err != nil {
		return nil, err
	}
	return ret, nil
}

func (b *builder) SortBytes(data [][]byte) {
	sort.Slice(data, func(i, j int) bool {
		return bytes.Compare(data[i], data[j]) < 0
	})
}

func (b *builder) start(da *DoubleArray, data []word.Word) error {
	da.init(0)
	b.id = 0
	return walker.Walk(data, func(prefix word.Word, branch []word.Code) error {
		return b.insert(da, prefix, branch)
	})
}

func (b *builder) insert(da *DoubleArray, prefix word.Word, branch []word.Code) error {
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
	offset := da.findValidOffset(branch)

	maxIndex := offset.Forward(branch[len(branch)-1])
	da.ensure(maxIndex)

	// nodes[at]にbranchを格納できるoffsetを指定
	da.nodes[index].SetOffset(offset)

	prev := word.NONE
	for _, c := range branch {
		// branch には同じラベルの枝が複数あることがある
		// もう追加していたらスキップ
		if c == prev {
			continue
		}
		next := offset.Forward(c)
		da.popNode(next)
		da.nodes[next].SetParent(index)

		if c == word.EOS {
			//終端マーク
			da.nodes[index].Terminate()
			da.nodes[next].SetOffset(b.id)
			b.id++
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
