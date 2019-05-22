package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() (*Server, *httptest.Server) {
	server := NewServer()
	testServer := httptest.NewServer(server.router)
	return server, testServer
}

func teardown(s *httptest.Server) {
	s.Close()
}

func TestCreateBasket(t *testing.T) {
	server, testServer := setup()
	defer teardown(testServer)

	r := createBasket(t, testServer)
	id := r["id"].(string)

	assert.NotNil(t, server.workers[id])
	assert.Regexp(t, "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$", id)
}

func TestAddItemNoBasket(t *testing.T) {
	_, testServer := setup()
	defer teardown(testServer)

	resp := addItemToBasket(t, testServer, "100", "VOUCHER")
	bs := readBody(t, resp)

	r := map[string]interface{}{}
	json.Unmarshal(bs, &r)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "basket not found", r["error"])
}

func TestAddItemHavingBasket(t *testing.T) {
	_, testServer := setup()
	defer teardown(testServer)

	r := createBasket(t, testServer)
	addItemToBasket(t, testServer, r["id"].(string), "VOUCHER")
	addItemToBasket(t, testServer, r["id"].(string), "TSHIRT")
	addItemToBasket(t, testServer, r["id"].(string), "MUG")
	resp := addItemToBasket(t, testServer, r["id"].(string), "MUG")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	message := map[string]interface{}{}
	if err := json.Unmarshal(readBody(t, resp), &message); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 4000.0, message["total"])
}

func TestCloseBasket(t *testing.T) {
	server, testServer := setup()
	defer teardown(testServer)

	jsonBasket := createBasket(t, testServer)
	id := jsonBasket["id"].(string)
	_, err := server.getWorker(id)
	if err != nil {
		t.Fatal(err)
	}

	closeBasket(t, testServer, id)
	worker, err := server.getWorker(id)
	assert.Nil(t, worker)
	assert.IsType(t, &BasketNotFoundError{}, err)
}

func TestGetItems(t *testing.T) {
	server, testServer := setup()
	defer teardown(testServer)

	server.items = map[Code]Item{
		Code("VOUCHER"): Item{
			Code:  Code("VOURHCER"),
			Name:  "a voucher",
			Price: 250,
		},
		Code("LONGCLAW"): Item{
			Code:  Code("LONGCLAW"),
			Name:  "longclaw",
			Price: 1000000,
		},
	}

	resp, err := testServer.Client().Get(testServer.URL + "/items")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bs := readBody(t, resp)
	json := string(bs)

	assert.Equal(t, json, ""+
		"["+
		`{"code":"VOURHCER","name":"a voucher","price":250},`+
		`{"code":"LONGCLAW","name":"longclaw","price":1000000}`+
		"]")
}

func createBasket(t *testing.T, testServer *httptest.Server) map[string]interface{} {
	resp, err := testServer.Client().Post(testServer.URL+"/baskets", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	bs := readBody(t, resp)

	r := map[string]interface{}{}
	if err := json.Unmarshal(bs, &r); err != nil {
		t.Fatalf("%s, %v", bs, err)
	}

	return r
}

func closeBasket(t *testing.T, testServer *httptest.Server, basketID string) {
	req, err := http.NewRequest("PUT", testServer.URL+"/baskets/"+basketID, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func addItemToBasket(t *testing.T, testServer *httptest.Server, basketID, code string) *http.Response {
	req, err := http.NewRequest("PUT",
		testServer.URL+"/baskets/"+basketID+"/items",
		strings.NewReader(fmt.Sprintf(`{"code": "%s"}`, code)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("content-type", "application/json")

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	return resp
}

func readBody(t *testing.T, resp *http.Response) []byte {
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return bs
}
