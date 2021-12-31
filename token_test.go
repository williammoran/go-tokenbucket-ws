package main

import (
	"testing"
	"time"
)

func TestReplenish(t *testing.T) {
	b := makeTestBucket()
	k := b.Get("meow")
	k.lastReplenish = time.Time{}
	k.Replenish(b, time.Time{}.Add(5*time.Second))
	if k.current != 100000 {
		t.Logf("Expected 100000 but got %d", k.current)
		t.Fail()
	}
	k.current = 0
	k.lastReplenish = time.Time{}
	k.Replenish(b, time.Time{}.Add(5*time.Second))
	if k.current != 5 {
		t.Logf("Expected 5 but got %d", k.current)
		t.Fail()
	}
}
