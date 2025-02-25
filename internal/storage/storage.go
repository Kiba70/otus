package storage

import (
	"log/slog"
	"sync"
)

const (
	defaultCountSeconds = 86400
)

type (
	Storage struct {
		elem         []any
		headPoint    int
		countSeconds int

		sync.RWMutex
	}
)

func New(seconds ...int) *Storage {
	countSeconds := defaultCountSeconds // Размер буфера по умолчанию

	switch {
	case len(seconds) > 1:
		return nil // Не коррекный параметр
	case len(seconds) == 1:
		countSeconds = seconds[0]
	}

	result := &Storage{
		elem:         make([]any, 0, countSeconds),
		headPoint:    0,
		countSeconds: countSeconds,
	}

	return result
}

func (s *Storage) Add(elem any) {
	s.Lock()
	defer s.Unlock()

	if len(s.elem) < s.countSeconds { // Буфер заполнен не весь - просто добавляем
		s.elem = append(s.elem, elem)
	} else {
		s.elem[s.headPoint] = elem
	}

	// Сдвигаем указатель на новую вставку на 1
	s.headPoint++
	if s.headPoint == s.countSeconds {
		s.headPoint = 0
	}
}

func (s *Storage) Get(m int) []any {
	s.RLock()
	defer s.RUnlock()

	if len(s.elem) < m {
		return nil // Буфер с данными не заполнен
	}

	result := make([]any, m)

	if s.headPoint-m >= 0 { // Попадаем в 1 слайс
		ncopy := copy(result, s.elem[s.headPoint-m:])
		slog.Debug("Storage copy", "num copyed", ncopy)
		return result
	}

	// Необходимые данные разделены на 2 части (начало и конец буфера)
	ncopy := copy(result, s.elem[:s.headPoint])
	slog.Debug("Storage copy part1", "num copyed", ncopy)
	ncopy = copy(result[s.headPoint:], s.elem[m-s.headPoint:])
	slog.Debug("Storage copy part1", "num copyed", ncopy)

	return result
}
