package doublearray

import (
	"fmt"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

func ForEach(da Nodes, yield func(word.Word, uint32) error) error {
	for i := node.Index(0); ; i++ {
		nod, err := da.At(i)
		if err != nil {
			return nil
		}
		if !nod.IsTerminal() {
			continue
		}
		key, err := getKey(da, nod.GetParent(), i)
		if err != nil {
			return err
		}
		dat, err := da.At(nod.GetChild(word.EOS))
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
func getKey(da Nodes, parent, child node.Index) (word.Word, error) {
	buf := make(word.Word, 0, 8)
	nod, err := da.At(parent)
	if err != nil {
		return nil, err
	}
	for {
		label := node.Label(nod.GetOffset(), child)
		if nod.GetChild(label) != child {
			return nil, fmt.Errorf("bad label(%d) want %d got %d",
				label,
				child,
				nod.GetChild(label),
			)
		}
		buf = append(buf, label)

		if !nod.HasParent() {
			break
		}

		next := nod.GetParent()
		nod, err = da.At(next)
		if err != nil {
			return nil, err
		}
		parent, child = next, parent
	}
	reverse(buf)
	return buf, nil
}

func reverse(w []word.Code) {
	i, j := 0, len(w)-1
	for i < j {
		w[i], w[j] = w[j], w[i]
		i++
		j--
	}
}
