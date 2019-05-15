package server

import (
	"fmt"
	"net/http"
)

func AddItemEndpoint(s *Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		worker, err := s.getWorker(r)
		if handleError(err, w) {
			return
		}

		codeStr := r.URL.Query().Get("code")
		if codeStr == "" {
			handleError(&MissingItemCodeError{}, w)
			return
		}
		code := Code(codeStr)
		getTotalChan := getTotalChanPool.Get().(chan Total)
		errorChan := errorChanPool.Get().(chan error)

		worker.AddItem <- AddItemAction{
			Code:         code,
			GetTotalChan: getTotalChan,
			ErrorChan:    errorChan,
		}

		select {
		case total := <-getTotalChan:
			message := fmt.Sprintf(`{"total": "%f"}`, total)
			w.Write([]byte(message))
			w.WriteHeader(200)
		case err := <-errorChan:
			handleError(err, w)
		}
	}
}
