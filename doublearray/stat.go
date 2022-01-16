package doublearray

import (
	"fmt"

	"github.com/ajiyoshi-vg/hairetsu/node"
)

type Stat struct {
	Size  int
	Node  int
	Empty int
	Leaf  int
}

func GetStat(da Nodes) Stat {
	ret := Stat{}
	for i := 0; ; i++ {
		x, err := da.At(node.Index(i))
		if err != nil {
			return ret
		}
		ret.Size++
		if x.IsTerminal() {
			ret.Leaf++
		}
		if x.IsUsed() {
			ret.Node++
		} else {
			ret.Empty++
		}
	}
}

func (s Stat) String() string {
	filled := float64(s.Node) / float64(s.Size)
	return fmt.Sprintf(`{"size":%d, "node":%d, "empty":%d, "leaf":%d, "filled":%f)`,
		s.Size,
		s.Node,
		s.Empty,
		s.Leaf,
		filled,
	)
}
