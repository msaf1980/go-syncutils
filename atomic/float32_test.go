// Copyright (c) 2020 Uber Technologies, Inc.
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFloat32(t *testing.T) {
	atom := NewFloat32(4.2)

	require.Equal(t, float32(4.2), atom.Load(), "Load didn't work.")

	require.True(t, atom.CompareAndSwap(4.2, 0.5), "CAS didn't report a swap.")
	require.Equal(t, float32(0.5), atom.Load(), "CAS didn't set the correct value.")
	require.False(t, atom.CompareAndSwap(0.0, 1.5), "CAS reported a swap.")

	atom.Store(42.0)
	require.Equal(t, float32(42.0), atom.Load(), "Store didn't set the correct value.")
	require.Equal(t, float32(42.5), atom.Add(0.5), "Add didn't work.")
	require.Equal(t, float32(42.0), atom.Sub(0.5), "Sub didn't work.")

	require.Equal(t, float32(42.0), atom.Swap(45.0), "Swap didn't return the old value.")
	require.Equal(t, float32(45.0), atom.Load(), "Swap didn't set the correct value.")

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, "42.5", NewFloat32(42.5).String(),
			"String() returned an unexpected value.")
	})
}
