package main

import (
	"time"
)

type token struct {
	//lock          sync.Mutex
	current       int64
	lastReplenish time.Time
}

func (t *token) Replenish(b *bucket, to time.Time) {
	amount := to.Sub(t.lastReplenish) / b.replenish
	t.current += int64(amount)
	if t.current > b.max {
		t.current = b.max
	}
	t.lastReplenish = to
}
