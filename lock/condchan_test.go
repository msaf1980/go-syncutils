package lock

import (
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/msaf1980/go-syncutils/atomic"
)

func TestCondChanSelectBroadcast(t *testing.T) {
	var cc CondChan

	startChan := make(chan struct{})
	finishChan := make(chan struct{})
	const waiterCnt = 100
	for i := 0; i < waiterCnt; i++ {
		go func() {
			startChan <- struct{}{}
			cc.L.Lock()
			cc.Wait()
			cc.L.Unlock()
			finishChan <- struct{}{}
		}()
	}

	for i := 0; i < waiterCnt; i++ {
		<-startChan
	}

	time.Sleep(time.Millisecond * 10)
	runtime.Gosched()
	cc.L.Lock()
	cc.Broadcast()
	cc.L.Unlock()

	for i := 0; i < waiterCnt; i++ {
		<-finishChan
	}
}

func TestCondChanSignalBroadcastRace(t *testing.T) {
	var cc CondChan
	finishChan := make(chan struct{})
	const times = 1000000

	go func() {
		for i := 0; i < times; i++ {
			cc.L.Lock()
			cc.Signal()
			cc.L.Unlock()
		}
		finishChan <- struct{}{}
	}()

	go func() {
		for i := 0; i < times; i++ {
			cc.L.Lock()
			cc.Broadcast()
			cc.L.Unlock()
		}
		finishChan <- struct{}{}
	}()

	<-finishChan
	<-finishChan
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tests below are taken from sync.Cond package ////////////////////////////////////////////////////////////////////////

func TestCondChanSignal(t *testing.T) {
	var cc CondChan
	n := 2
	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			running <- true
			cc.L.Lock()
			cc.Wait()
			cc.L.Unlock()
			awake <- true
		}()
	}
	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}
	time.Sleep(time.Millisecond * 10)
	for n > 0 {
		cc.L.Lock()
		cc.Signal()
		cc.L.Unlock()
		<-awake // Will deadlock if no goroutine wakes up
		n--
	}
}

func TestCondChanSignalGenerations(t *testing.T) {
	var cc CondChan
	n := 100
	running := make(chan bool, n)
	awake := make(chan int, n)
	for i := 0; i < n; i++ {
		go func(i int) {
			running <- true
			cc.L.Lock()
			cc.Wait()
			cc.L.Unlock()
			awake <- i
		}(i)
		if i > 0 {
			a := <-awake
			if a != i-1 {
				t.Fatalf("wrong goroutine woke up: want %d, got %d", i-1, a)
			}
		}
		<-running
		cc.L.Lock()
		cc.Signal()
		cc.L.Unlock()
		time.Sleep(time.Millisecond * 10)
		cc.L.Lock()
		cc.Signal()
		cc.L.Unlock()
	}
}

func TestCondChanSignalStealing(t *testing.T) {
	for iters := 0; iters < 1000; iters++ {
		var cc CondChan

		// Start a waiter.
		ch := make(chan struct{})
		go func() {
			ch <- struct{}{}
			cc.L.Lock()
			cc.Wait()
			cc.L.Unlock()
			ch <- struct{}{}
		}()

		<-ch

		time.Sleep(time.Millisecond * 20)

		// We know that the waiter is in the cond.Wait() call because we
		// synchronized with it, then acquired/released the mutex it was
		// holding when we synchronized.
		//
		// Start two goroutines that will race: one will broadcast on
		// the cond var, the other will wait on it.
		//
		// The new waiter may or may not get notified, but the first one
		// has to be notified.
		run := atomic.NewBool(true)
		go func() {
			cc.L.Lock()
			cc.Broadcast()
			cc.L.Unlock()
		}()

		go func() {
			for run.Load() {
				cc.L.Lock()
				cc.Wait()
				cc.L.Unlock()
			}
		}()

		// Check that the first waiter does get signaled.
		select {
		case <-ch:
		case <-time.After(2 * time.Second):
			t.Fatalf("First waiter didn't get broadcast.")
		}

		// Release the second waiter in case it didn't get the
		// broadcast.
		run.Store(false)
		cc.L.Lock()
		cc.Broadcast()
		cc.L.Unlock()
	}
}

func TestChanCondSignalRace(t *testing.T) {

	var (
		wgStart sync.WaitGroup
		wg      sync.WaitGroup
		cc      CondChan
	)

	n := 100
	clients := 200

	for i := 0; i < n; i++ {
		wgStart.Add(2 * clients)
		wg.Add(2 * clients)
		for j := 0; j < clients; j++ {
			for k := 0; k < 2; k++ {
				go func() {
					wgStart.Done()
					cc.L.Lock()
					cc.Wait()
					cc.L.Unlock()
					wg.Done()
				}()
			}
		}

		wgStart.Wait()
		runtime.Gosched()
		time.Sleep(time.Millisecond * 10)
		for j := 0; j < clients; j++ {
			runtime.Gosched()
			cc.L.Lock()
			cc.Signal()
			cc.Signal()
			cc.L.Unlock()
			runtime.Gosched()
		}
		wg.Wait()
	}

}

func TestCondChanSignalRace3(t *testing.T) {
	var (
		x  atomic.Int64
		cc CondChan
	)
	done := make(chan bool)
	step1 := make(chan bool)
	step2 := make(chan bool)

	go func() {
		x.Store(1)
		cc.L.Lock()
		step1 <- true
		cc.Wait()
		cc.L.Unlock()
		<-step2
		if x.Load() != 2 {
			cc.L.Unlock()
			t.Error("want 2")
		}
		x.Store(3)
		time.Sleep(time.Millisecond)
		cc.L.Lock()
		cc.Signal()
		cc.L.Unlock()
		done <- true
	}()
	go func() {
		for {
			if x.Load() == 1 {
				<-step1
				x.Store(2)
				// wait for lock and wait in third goroutine
				cc.L.Lock()
				cc.Signal()
				cc.L.Unlock()
				step2 <- true
				break
			}
			runtime.Gosched()
		}
		done <- true
	}()
	go func() {
		for {
			if x.Load() == 2 {
				cc.L.Lock()
				cc.Wait()
				cc.L.Unlock()
				if x.Load() != 3 {
					t.Error("want 3")
				}
				break
			}
			if x.Load() == 3 {
				break
			}
			runtime.Gosched()
		}
		done <- true
	}()
	<-done
	<-done
	<-done
}

func TestChanCondBroadcastRace(t *testing.T) {

	var (
		wgStart sync.WaitGroup
		wg      sync.WaitGroup
		cc      CondChan
	)

	n := 100
	clients := 200

	for i := 0; i < n; i++ {
		wgStart.Add(2 * clients)
		wg.Add(2 * clients)
		for j := 0; j < clients; j++ {
			for k := 0; k < 2; k++ {
				go func() {
					wgStart.Done()
					cc.L.Lock()
					cc.Wait()
					cc.L.Unlock()
					wg.Done()
				}()
			}
		}

		wgStart.Wait()
		time.Sleep(time.Millisecond * 20)
		runtime.Gosched()
		cc.L.Lock()
		cc.Broadcast()
		cc.L.Unlock()
		wg.Wait()
	}

}

func benchmarkCondChanBroadcast(b *testing.B, waiters int) {
	var (
		cc      CondChan
		wgStart sync.WaitGroup
	)

	done := make(chan bool)
	id := 0

	wgStart.Add(waiters + 1)
	for routine := 0; routine < waiters+1; routine++ {
		go func() {
			wgStart.Done()
			wgStart.Wait()
			for i := 0; i < b.N; i++ {
				cc.L.Lock()
				if id == -1 {
					cc.L.Unlock()
					break
				}
				id++
				if id == waiters+1 {
					id = 0
					cc.Broadcast()
				} else {
					cc.Wait()
				}
				cc.L.Unlock()
			}
			cc.L.Lock()
			id = -1
			cc.Broadcast()
			cc.L.Unlock()
			done <- true
		}()
	}
	for routine := 0; routine < waiters+1; routine++ {
		<-done
	}

}

func BenchmarkCondChanBroadcast(b *testing.B) {
	tests := []int{1, 5, 10, 20, 50, 100, 1000}
	for _, waiters := range tests {
		b.Run(strconv.Itoa(waiters), func(b *testing.B) {
			benchmarkCondChanBroadcast(b, waiters)
		})
	}
}

func benchmarkCondChanSignal(b *testing.B, waiters int) {
	var (
		cc      CondChan
		wgStart sync.WaitGroup
	)

	run := atomic.NewBool(true)

	done := make(chan bool)
	// awake := make(chan struct{}, waiters*10)

	wgStart.Add(waiters)
	for routine := 0; routine < waiters; routine++ {
		go func() {
			wgStart.Done()
			for run.Load() {
				runtime.Gosched()
				cc.L.Lock()
				cc.Wait()
				cc.L.Unlock()
			}
		}()
	}
	wgStart.Wait()

	wgStart.Add(waiters)
	for routine := 0; routine < waiters; routine++ {
		go func() {
			wgStart.Done()
			wgStart.Wait()
			for i := 0; i < b.N; i++ {
				cc.L.Lock()
				cc.Signal()
				cc.L.Unlock()
			}
			done <- true
			cc.L.Lock()
			cc.Broadcast()
			cc.L.Unlock()
		}()
	}
	for routine := 0; routine < waiters; routine++ {
		<-done
	}
	// close(awake)
	run.Store(false)
}

func BenchmarkCondChanSignal(b *testing.B) {
	tests := []int{1, 5, 10, 20, 50, 100, 1000}
	for _, waiters := range tests {
		b.Run(strconv.Itoa(waiters), func(b *testing.B) {
			benchmarkCondChanSignal(b, waiters)
		})
	}
}
