package web

import (
	"context"
	"net/http"
	"otus/internal/config"
	"sync"
)

func Start(ctx context.Context, wgGlobal *sync.WaitGroup) error {
	s := &http.Server{
		Addr: ":" + config.Port,
		// Handler: myHandler,

	}

	go func() {
		defer wgGlobal.Done()
		s.ListenAndServe()
	}()

	return nil
}
