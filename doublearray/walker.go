package doublearray

import (
	"fmt"
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

// WalkLeaf : Enumerate the leaf nodes.
// This is for testing purposes, so it does not work efficiently.
func WalkLeaf(da Nodes, yield func(word.Word, uint32) error) error {
	for i := node.Index(0); ; i++ {
		target, err := da.At(i)
		if err != nil {
			return nil
		}
		if !target.IsTerminal() {
			continue
		}
		key, err := getKey(da, target.GetParent(), i)
		if err != nil {
			return err
		}
		dat, err := da.At(target.GetChild(word.EOS))
		if err != nil {
			return err
		}
		if dat.GetParent() == i {
			if err := yield(key, uint32(dat.GetOffset())); err != nil {
				return err
			}
		}
	}
}
func Leafs(da Nodes) iter.Seq[item.Item] {
	return func(yield func(item.Item) bool) {
		_ = WalkLeaf(da, func(key word.Word, value uint32) error {
			if !yield(item.New(key, value)) {
				return fmt.Errorf("iteration stop")
			}
			return nil
		})
	}
}
func getKey(da Nodes, parent, child node.Index) (word.Word, error) {
	buf := make(word.Word, 0, 8)
	target, err := da.At(parent)
	if err != nil {
		return nil, err
	}
	for {
		label := node.Label(target.GetOffset(), child)
		if target.GetChild(label) != child {
			return nil, fmt.Errorf("bad label(%d) want %d got %d",
				label,
				child,
				target.GetChild(label),
			)
		}
		buf = append(buf, label)

		if !target.HasParent() {
			break
		}

		next := target.GetParent()
		target, err = da.At(next)
		if err != nil {
			return nil, err
		}
		parent, child = next, parent
	}
	word.Reverse(buf)
	return buf, nil
}
