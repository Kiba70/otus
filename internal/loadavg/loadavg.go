package loadavg

import "otus/internal/storage"

const (
	fileName = "/proc/loadavg"
)

type (
	statistic struct {
		one     float32
		five    float32
		fifteen float32
	}
)

func Start() error {
	s := storage.New()
	_ = s
	return nil
}
