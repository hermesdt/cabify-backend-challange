package server

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	workers map[string]*Worker
	router  *chi.Mux
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
	router := chi.NewRouter()
	s := &Server{
		workers: make(map[string]*Worker),
		router:  router,
	}

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/baskets", func(r chi.Router) {
		r.Post("/", CreateBasket(s))
		r.Route("/{id}", func(r chi.Router) {
			r.Use(workerCtx(s))

			r.Put("/", CloseBasketEndpoint(s))
			r.Put("/items", AddItemEndpoint(s))
		})
	})

	return s
}

func (s *Server) Start() {
	http.ListenAndServe(":3000", s.router)
}

func workerCtx(s *Server) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			worker, err := s.getWorker(id)
			if handleError(err, w) {
				return
			}

			ctx := context.WithValue(r.Context(), "worker", worker)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (s *Server) getWorker(id string) (*Worker, error) {
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
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	bs, _ := json.Marshal(map[string]interface{}{"error": err.Error()})
	w.Write(bs)

	return true
}
