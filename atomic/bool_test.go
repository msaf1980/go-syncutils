package atomic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	atom := NewBool(false)
	require.False(t, atom.Toggle(), "Expected Toggle to return previous value.")
	require.True(t, atom.Toggle(), "Expected Toggle to return previous value.")
	require.False(t, atom.Toggle(), "Expected Toggle to return previous value.")
	require.True(t, atom.Load(), "Unexpected state after swap.")

	require.True(t, atom.CompareAndSwap(true, true), "CAS should swap when old matches")
	require.True(t, atom.Load(), "CAS should have no effect")
	require.True(t, atom.CompareAndSwap(true, false), "CAS should swap when old matches")
	require.False(t, atom.Load(), "CAS should have modified the value")
	require.False(t, atom.CompareAndSwap(true, false), "CAS should fail on old mismatch")
	require.False(t, atom.Load(), "CAS should not have modified the value")

	atom.Store(false)
	require.False(t, atom.Load(), "Unexpected state after store.")

	prev := atom.Swap(false)
	require.False(t, prev, "Expected Swap to return previous value.")

	prev = atom.Swap(true)
	require.False(t, prev, "Expected Swap to return previous value.")

	t.Run("String", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			assert.Equal(t, "true", NewBool(true).String(),
				"String() returned an unexpected value.")
		})

		t.Run("false", func(t *testing.T) {
			var b Bool
			assert.Equal(t, "false", b.String(),
				"String() returned an unexpected value.")
		})
	})
}

func TestBool_InitializeDefaults(t *testing.T) {
	tests := []struct {
		msg     string
		newBool func() *Bool
	}{
		{
			msg: "Uninitialized",
			newBool: func() *Bool {
				var b Bool
				return &b
			},
		},
		{
			msg: "NewBool with default",
			newBool: func() *Bool {
				return NewBool(false)
			},
		},
		{
			msg: "Bool swapped with default",
			newBool: func() *Bool {
				b := NewBool(true)
				b.Swap(false)
				return b
			},
		},
		{
			msg: "Bool CAS'd with default",
			newBool: func() *Bool {
				b := NewBool(true)
				b.CompareAndSwap(true, false)
				return b
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {

			t.Run("String", func(t *testing.T) {
				b := tt.newBool()
				assert.Equal(t, "false", b.String())
			})

			t.Run("CompareAndSwap", func(t *testing.T) {
				b := tt.newBool()
				require.True(t, b.CompareAndSwap(false, true))
				assert.Equal(t, true, b.Load())
			})

			t.Run("Swap", func(t *testing.T) {
				b := tt.newBool()
				assert.Equal(t, false, b.Swap(true))
			})
		})
	}
}
