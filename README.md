# go-syncutils

**go-lock** is a Golang library implementing some missing synchronization primitives in golang standart library:
## mutex
- Mutex and RWMutex with timeout mechanism (with Trylock and Lock with timeout)
- PMutex is promote (RWMutex with promote from read lock to read write lock) with timeout mechanism (with Trylock and Lock with timeout)

## atomic
Simple wrappers for primitive types to enforce atomic access.

### Usage

The standard library's `sync/atomic` is powerful, but it's easy to forget which
variables must be accessed atomically. `go.uber.org/atomic` preserves all the
functionality of the standard library, but wraps the primitive types to
provide a safer, more convenient API.

```go
var atom atomic.Uint32
atom.Store(42)
atom.Sub(2)
atom.CompareAndSwap(40, 11)
```

## Installation

```sh
go get github.com/msaf1980/go-syncutils
```
