package storage_test

import (
	"fmt"
	"otus/internal/storage"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	s *storage.Storage
)

func TestStorage(t *testing.T) {
	t.Run("Creating storage", func(t *testing.T) {
		s = storage.New(10)
		require.NotNil(t, s)

		for i := 0; i < 15; i++ {
			s.Add(i)
		}

		g := s.Get(15)
		fmt.Println("G:", g[0].(int))
		require.NotNil(t, g)
		require.Equal(t, len(g), 15)
	})

	// t.Run("Adding data", func(t *testing.T) {
	// for i := 0; i < 15; i++ {
	// 	s.Add(i)
	// }
	// })

	// t.Run("Get 5 elements", func(t *testing.T) {
	// 	g := s.Get(5)
	// 	fmt.Println("G:", g)
	// 	require.Equal(t, len(g), 5)
	// })
}
