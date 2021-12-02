package main

import (
	"sync"
	"time"
)

func NewBucket(
	init, max int64, replenish time.Duration,
) *bucket {
	return &bucket{
		tokens:    make(map[string]*token),
		init:      init,
		max:       max,
		replenish: replenish,
	}
}

type use struct {
	Used      int64
	Remaining int64
}

type bucket struct {
	newLock   sync.Mutex
	tokens    map[string]*token
	init      int64
	max       int64
	replenish time.Duration
}

func (b *bucket) Use(id string, amount int64) use {
	b.newLock.Lock()
	defer b.newLock.Unlock()
	t := b.Get(id)
	//t.lock.Lock()
	//defer t.lock.Unlock()
	t.Replenish(b, time.Now())
	return b.use(t, amount)
}

func (b *bucket) use(t *token, amount int64) use {
	if t.current >= amount {
		t.current -= amount
	} else {
		amount = t.current
		t.current = 0
	}
	return use{
		Used:      amount,
		Remaining: t.current,
	}
}

func (b *bucket) Get(id string) *token {
	t, found := b.tokens[id]
	if !found {
		t = &token{
			current: b.init,
		}
		b.tokens[id] = t
	}
	return t
}
