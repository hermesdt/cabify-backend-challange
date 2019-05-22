package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func CreateBasket(s *Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		worker := NewWorker()
		s.workers[worker.GetId()] = worker
		worker.Start()

		bs, err := worker.Basket.MarshalJSON()
		if handleError(err, w) {
			return
		}

		writeJson(w, bs, http.StatusCreated)
	}
}

func AddItemEndpoint(s *Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		worker := r.Context().Value("worker").(*Worker)

		bs, err := ioutil.ReadAll(r.Body)
		if handleError(err, w) {
			return
		}

		codePayload := map[string]interface{}{}
		err = json.Unmarshal(bs, &codePayload)
		if handleError(err, w) {
			return
		}
		var code string
		var ok bool
		if code, ok = codePayload["code"].(string); !ok || code == "" {
			handleError(&MissingItemCodeError{}, w)
			return
		}

		item, ok := s.items[Code(code)]
		if !ok {
			handleError(&UnkownItemError{}, w)
			return
		}
		getBasket := make(chan Basket)
		errorChan := make(chan error)

		worker.AddItem <- AddItemAction{
			Item:      item,
			GetBasket: getBasket,
			ErrorChan: errorChan,
		}

		select {
		case basket := <-getBasket:
			bs, err := basket.MarshalJSON()
			if handleError(err, w) {
				return
			}
			writeJson(w, bs, http.StatusOK)
		case err := <-errorChan:
			handleError(err, w)
		}
	}
}

func CloseBasketEndpoint(s *Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		worker := r.Context().Value("worker").(*Worker)

		getBasket := make(chan Basket)

		worker.Close <- CloseAction{
			GetBasket: getBasket,
		}
		basket := <-getBasket
		delete(s.workers, worker.GetId())

		bs, err := basket.MarshalJSON()
		if handleError(err, w) {
			return
		}

		writeJson(w, bs, http.StatusOK)
	}
}

func GetItemsEndpoint(s *Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var items []Item
		for _, item := range s.items {
			items = append(items, item)
		}

		bs, err := json.Marshal(items)
		if handleError(err, w) {
			return
		}

		writeJson(w, bs, 200)
	}
}

func writeJson(w http.ResponseWriter, data []byte, status int) {
	w.WriteHeader(status)
	w.Header().Add("content-type", "application/json")

	w.Write(data)
}
