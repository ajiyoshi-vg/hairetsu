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
	if u.bar == nil {
		u.bar = uiprogress.AddBar(n).AppendCompleted().PrependElapsed()
	} else {
		u.bar.Total += n
	}
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
	if p.bar == nil {
		p.bar = progressbar.Default(int64(n))
	} else {
		p.bar.ChangeMax(n + p.bar.GetMax())
	}
}

func (p *ProgressBar) Add(n int) {
	if err := p.bar.Add(n); err != nil {
		slog.Debug("progress", "err", err)
	}
}
