package main

import (
	"testing"
	"time"
)

func TestReplenish(t *testing.T) {
	b := makeTestBucket()
	k := b.Get("meow")
	k.Replenish(b, time.Time{}.Add(5*time.Second))
	if k.current != 100 {
		t.Logf("Expected 100 but got %d", k.current)
		t.Fail()
	}
	k.current = 0
	k.Replenish(b, time.Time{}.Add(5*time.Second))
	if k.current != 5 {
		t.Logf("Expected 5 but got %d", k.current)
		t.Fail()
	}
}
