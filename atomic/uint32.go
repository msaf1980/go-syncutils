// @generated Code generated by gen-atomicint.

// Copyright (c) 2020-2023 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package atomic

import (
	"strconv"
	"sync/atomic"
)

// Uint32 is an atomic wrapper around uint32.
type Uint32 struct {
	_ noCopy

	v uint32
}

// NewUint32 creates a new Uint32.
func NewUint32(val uint32) *Uint32 {
	return &Uint32{v: val}
}

// Load atomically loads the wrapped value.
func (i *Uint32) Load() uint32 {
	return atomic.LoadUint32(&i.v)
}

// Add atomically adds to the wrapped uint32 and returns the new value.
func (i *Uint32) Add(delta uint32) uint32 {
	return atomic.AddUint32(&i.v, delta)
}

// Sub atomically subtracts from the wrapped uint32 and returns the new value.
func (i *Uint32) Sub(delta uint32) uint32 {
	return atomic.AddUint32(&i.v, ^(delta - 1))
}

// Inc atomically increments the wrapped uint32 and returns the new value.
func (i *Uint32) Inc() uint32 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped uint32 and returns the new value.
func (i *Uint32) Dec() uint32 {
	return i.Sub(1)
}

// CompareAndSwap is an atomic compare-and-swap.
func (i *Uint32) CompareAndSwap(old, new uint32) (swapped bool) {
	return atomic.CompareAndSwapUint32(&i.v, old, new)
}

// Store atomically stores the passed value.
func (i *Uint32) Store(val uint32) {
	atomic.StoreUint32(&i.v, val)
}

// Swap atomically swaps the wrapped uint32 and returns the old value.
func (i *Uint32) Swap(val uint32) (old uint32) {
	return atomic.SwapUint32(&i.v, val)
}

// String encodes the wrapped value as a string.
func (i *Uint32) String() string {
	v := i.Load()
	return strconv.FormatUint(uint64(v), 10)
}
