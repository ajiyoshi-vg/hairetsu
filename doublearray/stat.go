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
	ret := Stat{
		Size: da.length(),
	}
	for i := 0; i < da.length(); i++ {
		x, _ := da.at(node.Index(i))
		if x.IsTerminal() {
			ret.Leaf++
		}
		if x.IsUsed() {
			ret.Node++
		} else {
			ret.Empty++
		}
	}
	return ret
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
