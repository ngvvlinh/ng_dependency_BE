package lifecycle

import (
	"context"
	"net/http"
	"sync"

	"o.o/common/l"
)

var ll = l.New()

type HTTPServer struct {
	Name string
	*http.Server
}

func StartHTTP(notifyStop func(), servers ...HTTPServer) func() {
	var wg sync.WaitGroup
	wg.Add(len(servers))
	for _, s := range servers {
		// https://golang.org/doc/faq#closures_and_goroutines
		s := s
		go func() {
			defer wg.Done()
			defer func() { notifyStop() }()
			ll.S.Infof("HTTP server %v listening on %v", s.Name, s.Addr)
			err := s.ListenAndServe()
			switch err {
			case nil, http.ErrServerClosed:
				err = nil
			default:
				ll.S.Errorf("HTTP server %v error", s.Name, l.Error(err))
			}
		}()
	}

	return func() {
		for _, svr := range servers {
			go svr.Shutdown(context.Background())
		}
		wg.Wait()
	}
}
