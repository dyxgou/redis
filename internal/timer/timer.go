package timer

import (
	"strconv"
	"strings"
)

type timestamp struct {
	key  string
	time int64
}

func newTimestamp(key string, t int64) timestamp {
	return timestamp{
		key:  key,
		time: t,
	}
}

type Timer struct {
	ts []timestamp
	N  int
}

func New() *Timer {
	return &Timer{
		ts: make([]timestamp, 1024),
		N:  0,
	}
}

func (t *Timer) IsEmpty() bool {
	return t.N == 0
}

func (t *Timer) Insert(ts timestamp) {
	t.ts[t.N] = ts
	t.shiftUp(t.N)
	t.N++
}

func (t *Timer) Remove() timestamp {
	val := t.ts[0]
	t.swap(0, t.N-1)
	t.clearLast()
	t.shiftDown(0)

	return val
}

func (t *Timer) clearLast() {
	clear(t.ts[t.N-1 : t.N])
	t.N--
}

func (t *Timer) swap(i, j int) {
	t.ts[i], t.ts[j] = t.ts[j], t.ts[i]
}

func (t *Timer) less(i, j int) bool {
	return t.ts[i].time > t.ts[j].time
}

func (t *Timer) parent(i int) int {
	return (i - 1) >> 1
}

func (t *Timer) leftChild(i int) int {
	return 2*i + 1
}

func (t *Timer) rightChild(i int) int {
	return 2*i + 2
}

func (t *Timer) shiftUp(cur int) {
	for cur > 0 && t.less(t.parent(cur), cur) {
		parent := t.parent(cur)
		t.swap(parent, cur)
		cur = parent
	}
}

func (t *Timer) shiftDown(cur int) {
	for t.leftChild(cur) < t.N-1 {
		l := t.leftChild(cur)

		if l < t.N && t.less(l, l+1) {
			l++
		}

		if !t.less(cur, l) {
			break
		}

		t.swap(cur, l)
		cur = l
	}
}

func (t *Timer) String() string {
	var sb strings.Builder

	sb.WriteString("[ ")
	for i := 0; i < t.N; i++ {
		val := t.ts[i]
		sb.WriteString(val.key)
		sb.WriteByte('=')
		sb.WriteString(strconv.FormatInt(val.time, 10))

		if i != t.N-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(" ]")

	return sb.String()
}
