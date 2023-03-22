package mutex

import (
	"context"
	"fmt"
	"sync"
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

func BenchmarkDT_MutexLockUnlock(b *testing.B) {
	mx := sync.Mutex{}

	for i := 0; i < b.N; i++ {
		mx.Lock()

		go func() {
			mx.Unlock()
		}()

		mx.Lock()
		mx.Unlock()
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

// TODO: make normal test
