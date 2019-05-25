package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hermesdt/backend-challenge/pkg/model"
)

func (api *Api) CreateBasket(w http.ResponseWriter, r *http.Request) {
	basket := api.App.DB.BasketsStore.Create()

	bs, err := basket.MarshalJSON()
	if handleError(err, w) {
		return
	}

	writeJSON(w, bs, http.StatusCreated)
}

func (api *Api) CloseBasket(w http.ResponseWriter, r *http.Request) {
	basket := r.Context().Value("basket").(model.Basket)

	api.App.DB.BasketsStore.Delete(basket.ID.String())
	writeJSON(w, []byte("{}"), 200)
}

func (api *Api) AddItem(w http.ResponseWriter, r *http.Request) {
	basket := r.Context().Value("basket").(model.Basket)

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

	item, ok := model.ITEMS[model.Code(code)]
	if !ok {
		handleError(&UnkownItemError{}, w)
		return
	}

	basket.AddItem(item)
	err = api.App.DB.BasketsStore.Update(basket)
	if handleError(err, w) {
		return
	}

	bs, err = basket.MarshalJSON()
	if handleError(err, w) {
		return
	}

	writeJSON(w, bs, http.StatusOK)
}
