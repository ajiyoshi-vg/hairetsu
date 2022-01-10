package walker

import (
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Callback func(prefix word.Word, branch []word.Code) error

func Walk(data []word.Word, f Callback) error {
	return walkTrieNode(data, 0, len(data), 0, f)
}

func walkTrieNode(data []word.Word, begin, end, depth int, f Callback) error {
	// apply callback first
	if err := f(trieNode(data, begin, end, depth)); err != nil {
		return err
	}

	for begin < end {
		if data[begin].At(depth) != word.EOS {
			break
		}
		begin++
	}
	if begin == end {
		return nil
	}

	lastBegin := begin
	lastLabel := data[begin][depth]
	begin++
	for begin < end {
		label := data[begin][depth]
		if label != lastLabel {
			if err := walkTrieNode(data, lastBegin, begin, depth+1, f); err != nil {
				return err
			}
			lastBegin = begin
			lastLabel = label
		}
		begin++
	}
	return walkTrieNode(data, lastBegin, end, depth+1, f)
}

func trieNode(data []word.Word, begin, end, depth int) (word.Word, []word.Code) {
	prefix := data[begin][0:depth]
	branch := make([]word.Code, 0, end)
	for i := begin; i < end; i++ {
		branch = append(branch, data[i].At(depth))
	}
	return prefix, branch
}
