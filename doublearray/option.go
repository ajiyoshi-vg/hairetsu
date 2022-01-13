package doublearray

type Option func(*Builder)

func OptionProgress(p Progress) Option {
	return func(b *Builder) {
		b.progress = p
	}
}
