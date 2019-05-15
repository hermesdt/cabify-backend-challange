package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddItemNoBasket(t *testing.T) {
	s := NewServer()
	testServer := httptest.NewServer(s.mux)
	resp, err := testServer.Client().Get(testServer.URL + "/add_item?code=1")
	if err != nil {
		t.Fatal(err)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, `{"error": "basket not found"}`, string(bs))
	assert.Equal(t, 404, resp.StatusCode)
}

func TestAddItemHavingBasket(t *testing.T) {
	s := NewServer()
	testServer := httptest.NewServer(s.mux)
	resp, err := testServer.Client().Get(testServer.URL + "/new_basket")
	if err != nil {
		t.Fatal(err)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	message := &NewBasketResponse{}
	json.Unmarshal(bs, message)

	t.Log("message", message)
	t.FailNow()
}
