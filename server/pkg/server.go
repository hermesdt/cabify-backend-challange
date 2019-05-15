package server

import (
	"fmt"
	"net/http"
	"sync"
)

type Server struct {
	workers map[string]*Worker
	mux     *http.ServeMux
}

var getTotalChanPool = sync.Pool{
	New: func() interface{} {
		return make(chan Total)
	},
}

var errorChanPool = sync.Pool{
	New: func() interface{} {
		return make(chan error)
	},
}

func NewServer() *Server {
	mux := http.NewServeMux()
	s := &Server{
		mux: mux,
	}

	mux.HandleFunc("/add_item", AddItemEndpoint(s))

	return s
}

func (s *Server) Start() {
	http.ListenAndServe(":3000", s.mux)
}

func (s *Server) getWorker(r *http.Request) (*Worker, error) {
	id := r.Header.Get("X-Basket-ID")
	worker, ok := s.workers[id]
	if !ok {
		return nil, &BasketNotFoundError{}
	}

	return worker, nil
}

func handleError(err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}

	switch err.(type) {
	case *BasketNotFoundError:
		w.WriteHeader(404)
	default:
		w.WriteHeader(400)
	}

	message := fmt.Sprintf(`{"error": "%s"}`, err.Error())
	w.Write([]byte(message))

	return true
}
