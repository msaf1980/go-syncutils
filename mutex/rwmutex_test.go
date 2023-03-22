package mutex

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func BenchmarkRWMutexLockUnlock(b *testing.B) {
	mx := RWMutex{}

	for i := 0; i < b.N; i++ {
		mx.Lock()
		mx.Unlock()
	}
}

func BenchmarkRWMutexRLockRUnlock(b *testing.B) {
	mx := RWMutex{}

	for i := 0; i < b.N; i++ {
		mx.RLock()
		mx.RUnlock()
	}
}

func BenchmarkRWMutexLockUnlock_Parallel(b *testing.B) {
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

func BenchmarkRWMutexRLockRUnlock_Parallel(b *testing.B) {
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

func BenchmarkRWMutexRWLockRWUnlock_Parallel(b *testing.B) {
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

func BenchmarkRWMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := RWMutex{}

	for i := 0; i < b.N; i++ {
		mx.LockWithContext(ctx)
		mx.Unlock()
	}
}

func BenchmarkDT_RWMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := RWMutex{}

	for i := 0; i < b.N; i++ {
		mx.LockWithContext(ctx)

		go func() {
			mx.Unlock()
		}()

		mx.LockWithContext(ctx)
		mx.Unlock()
	}
}

func BenchmarkNT_RWMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := RWMutex{}

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

func BenchmarkN0T_RWMutexTryLockUnlock(b *testing.B) {
	ctx := context.Background()
	mx := RWMutex{}

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
func BenchmarkRWMutexTryRLockRUnlock(b *testing.B) {
	ctx := context.Background()
	mx := RWMutex{}

	for i := 0; i < b.N; i++ {
		mx.RLockWithContext(ctx)
		mx.RUnlock()
	}
}

func BenchmarkDT_RWMutexLockUnlock(b *testing.B) {
	mx := sync.RWMutex{}

	for i := 0; i < b.N; i++ {
		mx.Lock()

		go func() {
			mx.Unlock()
		}()

		mx.Lock()
		mx.Unlock()
	}
}

func BenchmarkNT_RWMutexLockUnlock(b *testing.B) {
	mx := sync.RWMutex{}

	k := 1000

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(k)

		mx.Lock()
		for j := 0; j < k; j++ {
			go func() {
				mx.Lock()

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

func BenchmarkN0T_RWMutexLockUnlock(b *testing.B) {
	mx := sync.RWMutex{}

	k := 1000

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(k)

		mx.Lock()
		for j := 0; j < k; j++ {
			go func() {
				mx.Lock()

				mx.Unlock()
				wg.Done()
			}()
		}
		mx.Unlock()

		wg.Wait()
	}
}

func BenchmarkMutexStdLockUnlock(b *testing.B) {
	mx := sync.Mutex{}

	for i := 0; i < b.N; i++ {
		mx.Lock()
		mx.Unlock()
	}
}

func BenchmarkDT_MutexStdLockUnlock(b *testing.B) {
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

func BenchmarkNT_MutexStdLockUnlock(b *testing.B) {
	mx := sync.Mutex{}

	k := 1000

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(k)

		mx.Lock()
		for j := 0; j < k; j++ {
			go func() {
				mx.Lock()

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

func BenchmarkN0T_MutexStdLockUnlock(b *testing.B) {
	mx := sync.Mutex{}

	k := 1000

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(k)

		mx.Lock()
		for j := 0; j < k; j++ {
			go func() {
				mx.Lock()

				mx.Unlock()
				wg.Done()
			}()
		}
		mx.Unlock()

		wg.Wait()
	}
}

func BenchmarkDT_N_MutexStdLockUnlock(b *testing.B) {
	mx := sync.Mutex{}

	for i := 0; i < b.N; i++ {
		mx.Lock()

		go func() {
		}()

		mx.Unlock()
	}
}

func BenchmarkNT_N_MutexStdLockUnlock(b *testing.B) {
	mx := sync.Mutex{}

	k := 1000

	for i := 0; i < b.N; i++ {
		mx.Lock()

		for j := 0; j < k; j++ {
			go func() {
			}()
		}

		mx.Unlock()
	}
}

func TestRWMutex(t *testing.T) {

	var mx RWMutex

	mx.RLock()
	if !mx.TryRLock() {
		t.Fatal("TestRWMutex TryRLock must success")
	}
	if mx.TryLock() {
		t.Fatal("TestRWMutex TryLock must fail")
	}

	mx.RUnlock()
	mx.RUnlock()

	mx.Lock()
	mx.Unlock()

	mx.Lock()
	t1 := mx.RLockWithDuration(time.Millisecond)
	if t1 {
		t.Fatal("TestRWMutex t1 fail R lock duration")
	}

	if mx.TryLock() {
		t.Fatal("TestRWMutex TryLock must fail")
	}

	if mx.TryRLock() {
		t.Fatal("TestRWMutex TryRLock must fail")
	}

	go func() {
		time.Sleep(5 * time.Millisecond)
		mx.Unlock()
	}()

	t2 := mx.RLockWithDuration(10 * time.Millisecond)
	t3 := mx.RLockWithDuration(10 * time.Millisecond)

	if !t2 {
		t.Fatal("TestRWMutex t2 fail R lock duration")
	}
	if !t3 {
		t.Fatal("TestRWMutex t3 fail R lock duration")
	}

}

// TODO: make normal test
