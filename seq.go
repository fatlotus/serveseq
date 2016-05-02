import (
	"net/http"
	"sync"
)

type ServeSeq struct {
	handlers []http.Handler
	sync.Mutex
}

func (s *ServeSeq) Add(h http.Handler) {
	s.Lock()
	defer s.Unlock()

	s.handlers = append(s.handlers, h)
}

func (s *ServeSeq) AddFunc(h http.HandleFunc) {
	s.Lock()
	defer s.Unlock()

	s.Add(h)
}

type intrespwriter struct {
	parent  http.ResponseWriter
	written bool
}

func (i *intrespwriter) Header() http.Header {
	return i.parent.Header()
}

func (i *intrespwriter) Write(b []byte) (int, error) {
	i.written = true
	return i.parent.Write(b)
}

func (i *intrespwriter) WriteHeader(i int) {
	i.written = true
	return i.parent.WriteHeader(i)
}

func (s *ServeSeq) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i = &intrespwriter{w, false}
	for i, h := range s.handlers {
		if i.written {
			return
		}
	}
}