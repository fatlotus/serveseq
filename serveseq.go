package serveseq

import (
	"net/http"
	"sync"
)

// A ServeSeq runs a list of handlers in order, exiting after the first one
// has responded to the incoming request.
type ServeSeq struct {
	handlers []http.Handler
	sync.Mutex
}

// Create a new ServeSeq object.
func New() *ServeSeq {
	return &ServeSeq{
		handlers: make([]http.Handler, 0),
	}
}

// Run this handler after all the others.
func (s *ServeSeq) Next(h http.Handler) {
	s.Lock()
	defer s.Unlock()

	s.handlers = append(s.handlers, h)
}

// Runs this other handler after all the others.
func (s *ServeSeq) NextFunc(h http.HandlerFunc) {
	s.Next(h)
}

type intrespwriter struct {
	parent                 http.ResponseWriter
	written, skip, canskip bool
}

func (i *intrespwriter) Header() http.Header {
	return i.parent.Header()
}

func (i *intrespwriter) Write(b []byte) (int, error) {
	if i.skip {
		return len(b), nil
	}

	i.written = true
	return i.parent.Write(b)
}

func (i *intrespwriter) WriteHeader(c int) {
	if i.skip {
		return
	} else if c == 404 && i.canskip && !i.written {
		i.skip = true
		return
	}

	i.written = true
	i.parent.WriteHeader(c)
}

// Run this ServeSeq as an http.Handler.
func (s *ServeSeq) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i := &intrespwriter{parent: w, written: false}
	for idx, h := range s.handlers {
		i.canskip = (idx != len(s.handlers)-1)
		i.skip = false
		h.ServeHTTP(i, r)
		if i.written {
			return
		}
	}
}
