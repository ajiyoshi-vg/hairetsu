package doublearray

import "fmt"

type Stat struct {
	Size  int
	Node  int
	Empty int
	Leaf  int
}

func newStat(da *DoubleArray) Stat {
	ret := Stat{
		Size: len(da.nodes),
	}
	for _, x := range da.nodes {
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
