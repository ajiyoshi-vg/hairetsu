package progress

import (
	"log/slog"

	"github.com/gosuri/uiprogress"
	"github.com/schollz/progressbar/v3"
)

type UIProgress struct {
	bar *uiprogress.Bar
}

func (u *UIProgress) SetMax(n int) {
	uiprogress.Start()
	u.bar = uiprogress.AddBar(n)
	u.bar.AppendCompleted()
	u.bar.PrependElapsed()
}

func (u *UIProgress) Add(n int) {
	for range n {
		u.bar.Incr()
	}
}

type ProgressBar struct {
	bar *progressbar.ProgressBar
}

func (p *ProgressBar) SetMax(n int) {
	p.bar = progressbar.Default(int64(n))
}

func (p *ProgressBar) Add(n int) {
	if err := p.bar.Add(n); err != nil {
		slog.Debug("progress", "err", err)
	}
}
