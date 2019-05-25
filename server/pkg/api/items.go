package api

import (
	"encoding/json"
	"net/http"

	"github.com/hermesdt/backend-challenge/pkg/model"
)

func (api *Api) GetItems(w http.ResponseWriter, r *http.Request) {
	var items []model.Item
	for _, item := range api.App.DB.ItemsStore.GetAll() {
		items = append(items, item)
	}

	bs, err := json.Marshal(items)
	if handleError(err, w) {
		return
	}

	writeJSON(w, bs, 200)
}
