package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
)

func CreateBasket(s *Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		worker := NewWorker()
		s.workers[worker.GetId()] = worker
		worker.Start()

		bs, err := json.Marshal(map[string]interface{}{"id": worker.GetId()})
		if handleError(err, w) {
			return
		}
		w.WriteHeader(201)
		w.Write(bs)
	}
}

type CodePayload struct {
	Code
}

func AddItemEndpoint(s *Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		basketID := chi.URLParam(r, "id")
		worker, err := s.getWorker(basketID)
		if handleError(err, w) {
			return
		}

		bs, err := ioutil.ReadAll(r.Body)
		if handleError(err, w) {
			return
		}

		codePayload := CodePayload{}
		err = json.Unmarshal(bs, &codePayload)
		if handleError(err, w) {
			return
		}

		if codePayload.Code == "" {
			handleError(&MissingItemCodeError{}, w)
			return
		}
		code := Code(codePayload.Code)
		getTotalChan := getTotalChanPool.Get().(chan Total)
		errorChan := errorChanPool.Get().(chan error)

		worker.AddItem <- AddItemAction{
			Code:         code,
			GetTotalChan: getTotalChan,
			ErrorChan:    errorChan,
		}

		select {
		case total := <-getTotalChan:
			bs, err := json.Marshal(map[string]interface{}{"total": total})
			if handleError(err, w) {
				return
			}

			w.Write(bs)
			w.WriteHeader(200)
		case err := <-errorChan:
			handleError(err, w)
		}
	}
}
