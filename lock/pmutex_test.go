package lock

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func BenchmarkPMutexLockUnlock(b *testing.B) {
	mx := PMutex{}

	for i := 0; i < b.N; i++ {
		mx.Lock()
		mx.Unlock()
	}
}

func BenchmarkPMutexRLockRUnlock(b *testing.B) {
	mx := PMutex{}

	for i := 0; i < b.N; i++ {
		mx.RLock()
		mx.RUnlock()
	}
}

func BenchmarkPMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := PMutex{}

	for i := 0; i < b.N; i++ {
		mx.LockWithContext(ctx)
		mx.Unlock()
	}
}

func BenchmarkDT_PMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := PMutex{}

	for i := 0; i < b.N; i++ {
		mx.LockWithContext(ctx)

		go func() {
			mx.Unlock()
		}()

		mx.LockWithContext(ctx)
		mx.Unlock()
	}
}

func BenchmarkNT_PMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := PMutex{}

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

func BenchmarkN0T_PMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := PMutex{}

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
func BenchmarkPMutexTryRLockRUnlock(b *testing.B) {
	ctx := context.Background()
	mx := PMutex{}

	for i := 0; i < b.N; i++ {
		mx.RLockWithContext(ctx)
		mx.RUnlock()
	}
}

func BenchmarkPMutexLockUnlock_Parallel(b *testing.B) {
	concurrencyLevels := []int{5, 10, 20, 50, 100, 1000}
	for _, clients := range concurrencyLevels {

		b.Run(fmt.Sprintf("%d", clients), func(b *testing.B) {
			wgStart := sync.WaitGroup{}
			wgStart.Add(clients)
			wg := sync.WaitGroup{}

			mx := PMutex{}
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

func BenchmarkPMutexRLockRUnlock_Parallel(b *testing.B) {
	concurrencyLevels := []int{5, 10, 20, 50, 100, 1000}
	for _, clients := range concurrencyLevels {

		b.Run(fmt.Sprintf("%d", clients), func(b *testing.B) {
			wgStart := sync.WaitGroup{}
			wgStart.Add(clients)
			wg := sync.WaitGroup{}

			mx := PMutex{}
			b.ResetTimer()
			for i := 0; i < clients; i++ {
				wg.Add(1)
				go func() {
					wgStart.Done()
					wgStart.Wait()
					// Test routine
					for n := 0; n < b.N; n++ {
						mx.RLock()
						mx.RUnlock()
					}
					// End test routine
					wg.Done()
				}()

			}

			wg.Wait()
		})
	}
}

func BenchmarkPMutexPromote_Parallel(b *testing.B) {
	concurrencyLevels := []int{5, 10, 20, 50, 100, 1000}
	for _, clients := range concurrencyLevels {

		b.Run(fmt.Sprintf("%d", clients), func(b *testing.B) {
			wgStart := sync.WaitGroup{}
			wgStart.Add(clients)
			wg := sync.WaitGroup{}

			mx := PMutex{}
			b.ResetTimer()
			for i := 0; i < clients; i++ {
				wg.Add(1)
				go func(promote bool) {
					wgStart.Done()
					wgStart.Wait()
					// Test routine
					if promote {
						mx.RLock()
						mx.Promote()
						mx.Unlock()
					} else {
						for n := 0; n < b.N; n++ {
							mx.RLock()
							mx.RUnlock()
						}
					}
					// End test routine
					wg.Done()
				}(i == 0)

			}

			wg.Wait()
		})
	}
}

func BenchmarkPMutexRWLockRWUnlock_Parallel(b *testing.B) {
	concurrencyLevels := []int{5, 10, 20, 50, 100, 1000}
	for _, clients := range concurrencyLevels {

		b.Run(fmt.Sprintf("%d", clients), func(b *testing.B) {
			wgStart := sync.WaitGroup{}
			wgStart.Add(clients)
			wg := sync.WaitGroup{}

			mx := RWMutex{}
			b.ResetTimer()
			for i := 0; i < clients; i++ {
				wg.Add(1)
				go func(read bool) {
					wgStart.Done()
					wgStart.Wait()
					// Test routine
					if read {
						for n := 0; n < b.N; n++ {
							mx.RLock()
							mx.RUnlock()
						}
					} else {
						for n := 0; n < b.N; n++ {
							mx.Lock()
							mx.Unlock()
						}
					}
					// End test routine
					wg.Done()
				}(i%2 == 0)

			}

			wg.Wait()
		})
	}
}

func TestPMutex(t *testing.T) {

	var mx PMutex

	mx.Lock()
	mx.Unlock()

	mx.Lock()
	t1 := mx.RLockWithTimeout(time.Millisecond)
	if t1 {
		t.Fatal("TestPMutex t1 fail R lock duration")
	}

	go func() {
		time.Sleep(5 * time.Millisecond)
		mx.Unlock()
	}()

	t2 := mx.RLockWithTimeout(50 * time.Millisecond)
	t3 := mx.RLockWithTimeout(50 * time.Millisecond)

	if !t2 {
		t.Fatal("TestPMutex t2 fail R lock duration")
	}
	if !t3 {
		t.Fatal("TestPMutex t3 fail R lock duration")
	}

}

func TestPMutexPromote(t *testing.T) {

	var mx PMutex

	mx.RLock()
	mx.Promote()

	t1 := mx.TryPromote()
	if t1 {
		t.Fatal("TestPMutex t1 must fail TryPromote")
	}

	mx.Unlock()

	mx.Lock()
	mx.Reduce()
	mx.RUnlock()

	mx.RLock()
	mx.RLock()
	go func() {
		time.Sleep(5 * time.Millisecond)
		mx.RUnlock()
	}()

	t2 := mx.PromoteWithTimeout(10 * time.Millisecond)
	if !t2 {
		t.Fatal("TestPMutex t2 fail Promote duration")
	}

	mx.Unlock()

	mx.Lock()
	mx.Reduce()

	t3 := mx.TryPromote()
	if !t3 {
		t.Fatal("TestPMutex t3 fail TryPromote")
	}
	mx.Reduce()

	mx.RUnlock()

}

// TODO: make normal test
