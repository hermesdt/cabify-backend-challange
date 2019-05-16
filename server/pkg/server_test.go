package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

	resp := addItemToBasket(t, testServer, "100")
	bs := readBody(t, resp)

	r := map[string]interface{}{}
	json.Unmarshal(bs, &r)

	assert.Equal(t, 404, resp.StatusCode)
	assert.Equal(t, r["error"], "basket not found")
}

func TestAddItemHavingBasket(t *testing.T) {
	_, testServer := setup()
	defer teardown(testServer)

	r := createBasket(t, testServer)
	resp := addItemToBasket(t, testServer, r["id"].(string))

	bs := readBody(t, resp)
	fmt.Println(string(bs))

	message := map[string]interface{}{}
	if err := json.Unmarshal(bs, &message); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 5.0, message["total"].(float64))
}

func createBasket(t *testing.T, testServer *httptest.Server) map[string]interface{} {
	resp, err := testServer.Client().Post(testServer.URL+"/baskets", "application/json", bytes.NewReader([]byte{}))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp.StatusCode, 201)

	bs := readBody(t, resp)

	r := map[string]interface{}{}
	if err := json.Unmarshal(bs, &r); err != nil {
		t.Fatalf("%s, %v", bs, err)
	}

	return r
}

func addItemToBasket(t *testing.T, testServer *httptest.Server, basketID string) *http.Response {
	resp, err := testServer.Client().Post(
		testServer.URL+"/baskets/"+basketID+"/items",
		"application/json",
		bytes.NewBufferString(`{"code": "VOUCHER"}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func readBody(t *testing.T, resp *http.Response) []byte {
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return bs
}
