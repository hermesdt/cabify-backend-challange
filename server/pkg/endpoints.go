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

		writeJson(w, map[string]interface{}{"id": worker.GetId()}, http.StatusCreated)
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
		getTotalChan := make(chan Total)
		errorChan := make(chan error)

		worker.AddItem <- AddItemAction{
			Item:         item,
			GetTotalChan: getTotalChan,
			ErrorChan:    errorChan,
		}

		select {
		case total := <-getTotalChan:
			writeJson(w, map[string]interface{}{"total": total}, http.StatusOK)
		case err := <-errorChan:
			handleError(err, w)
		}
	}
}

func CloseBasketEndpoint(s *Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		worker := r.Context().Value("worker").(*Worker)

		getTotalChan := make(chan Total)

		worker.Close <- CloseAction{
			GetTotalChan: getTotalChan,
		}
		total := <-getTotalChan
		delete(s.workers, worker.GetId())

		writeJson(w, map[string]interface{}{"total": total}, http.StatusOK)
	}
}

func GetItemsEndpoint(s *Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonItems := []map[string]interface{}{}
		for _, item := range s.items {
			jsonItems = append(jsonItems, item.asJson())
		}
		data := map[string]interface{}{
			"data": jsonItems,
		}
		writeJson(w, data, 200)
	}
}

func writeJson(w http.ResponseWriter, data map[string]interface{}, status int) {
	w.WriteHeader(status)
	w.Header().Add("content-type", "application/json")

	bs, err := json.Marshal(data)
	if handleError(err, w) {
		return
	}

	w.Write(bs)
}
