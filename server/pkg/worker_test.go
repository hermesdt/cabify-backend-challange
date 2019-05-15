package server

import "testing"

func TestRunAddItemClose(t *testing.T) {
	worker := NewWorker()
	worker.Run()
}
