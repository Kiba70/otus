package storage_test

import (
	"otus/internal/storage"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	s *storage.Storage[int]
)

func TestStorage(t *testing.T) {
	t.Run("Creating storage", func(t *testing.T) {
		s = storage.New[int](10)
		require.NotNil(t, s)
	})

	t.Run("Adding data", func(t *testing.T) {
		for i := 0; i < 15; i++ {
			s.Add(i)
		}
	})

	t.Run("Get 10 elements with len buffer 10", func(t *testing.T) {
		g := s.Get(10)
		require.NotNil(t, g)
		require.Equal(t, 10, len(g))
	})

	t.Run("Get 15 elements with len buffer 10", func(t *testing.T) {
		g := s.Get(15)
		require.Nil(t, g)
		require.Equal(t, 0, len(g))
	})
}
