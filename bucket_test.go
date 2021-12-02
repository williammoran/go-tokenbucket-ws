package main

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestGetCreates(t *testing.T) {
	b := makeTestBucket()
	k := b.Get("meow")
	if k.current != b.init {
		t.Logf("size %d != %d", k.current, b.init)
		t.Fail()
	}
}

func TestUse(t *testing.T) {
	b := makeTestBucket()
	k := b.Get("meow")
	r := b.use(k, 50)
	if r.Remaining != 50 || r.Used != 50 {
		t.Logf("Wrong: %+v", r)
		t.Fail()
	}
	if k.current != 50 {
		t.Logf("Didn't decriment bucket: %+v", k)
		t.Fail()
	}
}

func BenchmarkConcurrencySmallFew(b *testing.B) {
	bucket := makeTestBucket()
	wg := sync.WaitGroup{}
	b.ResetTimer()
	for m := 0; m < 4; m++ {
		wg.Add(1)
		go func() {
			for n := 0; n < b.N; n++ {
				v := rand.Intn(50)
				bucket.Use(strconv.Itoa(v), 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkConcurrencySmallMany(b *testing.B) {
	bucket := makeTestBucket()
	wg := sync.WaitGroup{}
	b.ResetTimer()
	for m := 0; m < 4; m++ {
		wg.Add(1)
		go func() {
			for n := 0; n < b.N; n++ {
				v := rand.Intn(9999)
				bucket.Use(strconv.Itoa(v), 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkConcurrencySmallVeryMany(b *testing.B) {
	bucket := makeTestBucket()
	wg := sync.WaitGroup{}
	b.ResetTimer()
	for m := 0; m < 4; m++ {
		wg.Add(1)
		go func() {
			for n := 0; n < b.N; n++ {
				v := rand.Intn(999999)
				bucket.Use(strconv.Itoa(v), 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkConcurrency(b *testing.B) {
	bucket := makeTestBucket()
	wg := sync.WaitGroup{}
	b.ResetTimer()
	for m := 0; m < 50; m++ {
		wg.Add(1)
		go func() {
			for n := 0; n < b.N; n++ {
				v := rand.Intn(50)
				bucket.Use(strconv.Itoa(v), 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkConcurrencyLarge(b *testing.B) {
	bucket := makeTestBucket()
	wg := sync.WaitGroup{}
	b.ResetTimer()
	for m := 0; m < 50; m++ {
		wg.Add(1)
		go func() {
			for n := 0; n < b.N; n++ {
				v := rand.Intn(9999)
				bucket.Use(strconv.Itoa(v), 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkConcurrencyLarger(b *testing.B) {
	bucket := makeTestBucket()
	wg := sync.WaitGroup{}
	b.ResetTimer()
	for m := 0; m < 50; m++ {
		wg.Add(1)
		go func() {
			for n := 0; n < b.N; n++ {
				v := rand.Intn(99999)
				bucket.Use(strconv.Itoa(v), 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func makeTestBucket() *bucket {
	return NewBucket(100000, 100000, time.Second)
}
