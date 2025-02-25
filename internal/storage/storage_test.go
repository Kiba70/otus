package storage_test

import (
	"otus/internal/storage"
	"testing"

	"github.com/stretchr/testify/require"
)

// var (
// 	s *storage.Storage[int]
// )

type stype struct {
	s1 int
	s2 int
	s3 int
}

func TestStorageInt(t *testing.T) {
	var s *storage.Storage[int]
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

func TestStorageStruct(t *testing.T) {
	var s *storage.Storage[stype]
	t.Run("Creating storage", func(t *testing.T) {
		s = storage.New[stype](10)
		require.NotNil(t, s)
	})

	t.Run("Adding data", func(t *testing.T) {
		for i := 0; i < 15; i++ {
			s2 := &stype{
				s1: i,
				s2: i + 1,
				s3: i + 2,
			}
			s.Add(*s2)
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
