package lock

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkChanMutexLockUnlock(b *testing.B) {
	mx := NewChanMutex()

	for i := 0; i < b.N; i++ {
		mx.Lock()
		mx.Unlock()
	}
}

func BenchmarkChanMutexLockUnlock_Parallel(b *testing.B) {
	concurrencyLevels := []int{5, 10, 20, 50, 100, 1000}
	for _, clients := range concurrencyLevels {

		b.Run(fmt.Sprintf("%d", clients), func(b *testing.B) {
			wgStart := sync.WaitGroup{}
			wgStart.Add(clients)
			wg := sync.WaitGroup{}

			mx := NewChanMutex()
			b.ResetTimer()
			for i := 0; i < clients; i++ {
				wg.Add(1)
				go func() {
					wgStart.Done()
					wgStart.Wait()
					// Test routine
					for n := 0; n < b.N; n++ {
						mx.Lock()
						mx.Unlock()
					}
					// End test routine
					wg.Done()
				}()

			}

			wg.Wait()
		})
	}
}

func BenchmarkChanMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := NewChanMutex()

	for i := 0; i < b.N; i++ {
		mx.LockWithContext(ctx)
		mx.Unlock()
	}
}

func BenchmarkDT_ChanMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := NewChanMutex()

	for i := 0; i < b.N; i++ {
		mx.LockWithContext(ctx)

		go func() {
			mx.Unlock()
		}()

		mx.LockWithContext(ctx)
		mx.Unlock()
	}
}

func BenchmarkNT_ChanMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := NewChanMutex()

	k := 1000

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(k)

		mx.LockWithContext(ctx)
		for j := 0; j < k; j++ {
			go func() {
				mx.LockWithContext(ctx)

				go func() {
					mx.Unlock()
					wg.Done()
				}()
			}()
		}
		mx.Unlock()

		wg.Wait()
	}
}

func BenchmarkN0T_ChanMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := NewChanMutex()

	k := 1000

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(k)

		mx.LockWithContext(ctx)
		for j := 0; j < k; j++ {
			go func() {
				mx.LockWithContext(ctx)

				mx.Unlock()
				wg.Done()
			}()
		}
		mx.Unlock()

		wg.Wait()
	}
}

func TestChanMutex(t *testing.T) {

	mx := NewChanMutex()

	mx.Lock()
	mx.Unlock()

	mx.Lock()
	t1 := mx.LockWithTimeout(time.Millisecond)
	if t1 {
		t.Fatal("TestChanMutex t1 fail R lock duration")
	}

	if mx.TryLock() {
		t.Fatal("TestChanMutex TryLock must fail")
	}

	go func() {
		time.Sleep(5 * time.Millisecond)
		mx.Unlock()
	}()

	t2 := mx.LockWithTimeout(20 * time.Millisecond)

	if !t2 {
		t.Fatal("TestChanMutex t2 fail R lock duration")
	}

}

func TestChanMutex_Race(t *testing.T) {

	var (
		wgStart sync.WaitGroup
		wg      sync.WaitGroup
		locks   int64
	)

	mx := NewChanMutex()

	n := 1000
	clients := 1000
	wgStart.Add(clients)
	wg.Add(clients)

	// rw lock
	for j := 0; j < clients; j++ {
		go func() {
			wgStart.Done()
			wgStart.Wait()
			for i := 0; i < n; i++ {
				mx.Lock()
				time.Sleep(time.Nanosecond * 10)
				mx.Unlock()
				atomic.AddInt64(&locks, 1)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	want := int64(n * clients)
	if locks != want {
		t.Errorf("locks count = %d, want %d", locks, want)
	}
}

// TODO: make normal test
