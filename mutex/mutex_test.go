package mutex

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkMutexLockUnlock(b *testing.B) {
	mx := Mutex{}

	for i := 0; i < b.N; i++ {
		mx.Lock()
		mx.Unlock()
	}
}

func BenchmarkMutexLockUnlock_Parallel(b *testing.B) {
	concurrencyLevels := []int{5, 10, 20, 50, 100, 1000}
	for _, clients := range concurrencyLevels {

		b.Run(fmt.Sprintf("%d", clients), func(b *testing.B) {
			wgStart := sync.WaitGroup{}
			wgStart.Add(clients)
			wg := sync.WaitGroup{}

			mx := Mutex{}
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

func BenchmarkMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := Mutex{}

	for i := 0; i < b.N; i++ {
		mx.LockWithContext(ctx)
		mx.Unlock()
	}
}

func BenchmarkDT_MutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := Mutex{}

	for i := 0; i < b.N; i++ {
		mx.LockWithContext(ctx)

		go func() {
			mx.Unlock()
		}()

		mx.LockWithContext(ctx)
		mx.Unlock()
	}
}

func BenchmarkNT_MutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := Mutex{}

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

func BenchmarkN0T_MutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := Mutex{}

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

func TestMutex(t *testing.T) {

	var mx Mutex

	mx.Lock()
	mx.Unlock()

	mx.Lock()
	t1 := mx.LockWithDuration(time.Millisecond)
	if t1 {
		t.Fatal("TestMutex t1 fail R lock duration")
	}

	if mx.TryLock() {
		t.Fatal("TestMutex TryLock must fail")
	}

	go func() {
		time.Sleep(5 * time.Millisecond)
		mx.Unlock()
	}()

	t2 := mx.LockWithDuration(10 * time.Millisecond)

	if !t2 {
		t.Fatal("TestMutex t2 fail R lock duration")
	}

}

func TestMutex_Invalid_Unlock(t *testing.T) {
	var mx Mutex

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("TestRWMutex Unlock must panic")
		}
	}()

	mx.Unlock()
}

func TestMutex_Race(t *testing.T) {

	var (
		mx            RWMutex
		wgStart       sync.WaitGroup
		wg            sync.WaitGroup
		locks, rlocks int64
	)

	n := 1000
	clients := 1000
	wgStart.Add(2 * clients)
	wg.Add(2 * clients)

	// rw lock
	for j := 0; j < clients; j++ {
		go func() {
			wgStart.Done()
			wgStart.Wait()
			for i := 0; i < n; i++ {
				mx.Lock()
				time.Sleep(time.Microsecond)
				mx.Unlock()
				atomic.AddInt64(&locks, 1)
			}
			wg.Done()
		}()
	}

	// read lock
	for j := 0; j < clients; j++ {
		go func() {
			wgStart.Done()
			wgStart.Wait()
			for i := 0; i < n; i++ {
				mx.RLock()
				time.Sleep(time.Microsecond)
				mx.RUnlock()
				atomic.AddInt64(&rlocks, 1)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

// TODO: make normal test
