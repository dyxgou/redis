package timer

import (
	"context"
	"time"
)

type deleteFunc func(key string)

type Ticker struct {
	elapsed    int64
	ticker     *time.Ticker
	timer      *Timer
	deleteFunc func(key string)
}

func NewTicker() *Ticker {
	return &Ticker{
		elapsed: 0,
		timer:   New(1_000_000),
	}
}

func (t *Ticker) Insert(key string, ti int64) {
	t.timer.insert(NewTimestamp(key, ti+t.elapsed))
}

func (t *Ticker) Init(ctx context.Context, delFunc deleteFunc) {
	t.deleteFunc = delFunc
	t.ticker = time.NewTicker(time.Second)
	defer t.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.ticker.C:
			t.elapsed++
			t.handle()
		}
	}
}

func (t *Ticker) handle() {
	if t.timer.IsEmpty() {
		if t.elapsed != 0 {
			t.elapsed = 0
		}

		return
	}

	for ti := range t.timer.Exited(t.elapsed) {
		t.deleteFunc(ti.key)
	}
}
