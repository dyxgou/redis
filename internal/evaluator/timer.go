package evaluator

import (
	"time"
)

type timestamp struct {
	key  string
	time time.Time
}

func newTimestamp(k string, t uint64) *timestamp {

}

type Timer []timestamp

func (t Timer) Len() int { return len(t) }

func (t Timer) Less(i, j int) bool { return t[i].time.Before(t[j].time) }

func (t Timer) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t *Timer) Pop() any {
	n := len(*t)
	x := (*t)[n-1]
	clear((*t)[n-1 : n])

	return x
}
