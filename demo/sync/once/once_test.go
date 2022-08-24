package once

import (
	"fmt"
	"sync"
	"testing"
)

type one int

func (o *one) Increment() {
	*o++
}

func run(t *testing.T, once *sync.Once, o *one, c chan bool) {
	once.Do(func() { o.Increment() })
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestOnce(t *testing.T) {
	PrintOnce()
	PrintOnce()
	PrintOnce()
}

var once sync.Once

// 这个方法，不管调用几次，只会输出一次
func PrintOnce() {
	once.Do(func() {
		fmt.Println("只输出一次")
	})
}

func TestOncePanic(t *testing.T) {
	var once sync.Once
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("Once.Do did not panic")
			}
		}()
		once.Do(func() {
			panic("failed")
		})
	}()

	once.Do(func() {
		t.Fatalf("Once.Do called twice")
	})
}

func BenchmarkOnce(b *testing.B) {
	var once sync.Once
	f := func() {}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			once.Do(f)
		}
	})
}
