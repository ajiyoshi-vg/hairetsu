package doublearray

import "github.com/ajiyoshi-vg/external"

type Option func(*Builder)

func OptionProgress(p Progress) Option {
	return func(b *Builder) {
		b.progress = p
	}
}

func StreamChunkSize(n int) Option {
	return func(b *Builder) {
		b.sortOption = append(b.sortOption, external.ChunkSize(n))
	}
}

func Verbose(b *Builder) {
	b.verbose = true
}
