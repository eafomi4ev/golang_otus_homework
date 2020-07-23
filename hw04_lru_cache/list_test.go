package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("add and delete the only one existing item using link on Front", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)

		require.Equal(t, 1, l.Len())
		require.Same(t, l.Front(), l.Back())

		l.Remove(l.Front())

		require.Zero(t, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("add and delete the only one existing item using link on Back", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)

		require.Equal(t, 1, l.Len())
		require.Same(t, l.Front(), l.Back())

		l.Remove(l.Back())

		require.Zero(t, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("can pass string as item's value", func(t *testing.T) {
		l := NewList()
		testValue := "test value"

		l.PushFront(testValue)

		require.Equal(t, testValue, l.Front().Value)
	})

	t.Run("front and back have correct links after calling MoveToFront", func(t *testing.T) {
		l := NewList()

		l.PushFront(1) // 1
		l.PushFront(2) // 2 1
		l.PushFront(3) // 3 2 1

		l.MoveToFront(l.Front().Next.Next) // 1 3 2
		l.MoveToFront(l.Front().Next)      // 3 1 2

		require.Equal(t, 2, l.Back().Value)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
