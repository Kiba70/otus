////go:build unit

package storage_test

import (
	"testing"

	"github.com/Kiba70/otus/internal/storage"
	"github.com/stretchr/testify/require"
)

type stype struct {
	s1 int
	s2 int
	s3 int
}

func TestStorageInt(t *testing.T) {
	var s *storage.Storage[int]
	t.Run("Создаём storage размером 10 элементов", func(t *testing.T) {
		s = storage.New[int](10)
		require.NotNil(t, s)
	})

	t.Run("15 раз добавляем по элементу", func(_ *testing.T) {
		for i := range 15 {
			s.Add(i)
		}
	})

	t.Run("Get 10 elements with len buffer 10", func(t *testing.T) {
		g := s.Get(10)
		require.NotNil(t, g)
		require.Equal(t, 10, len(g))
		require.Equal(t, []int{10, 11, 12, 13, 14, 5, 6, 7, 8, 9}, g)
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

	t.Run("Adding data", func(_ *testing.T) {
		for i := range 15 {
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
