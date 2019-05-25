package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/hermesdt/backend-challenge/pkg/db"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/hermesdt/backend-challenge/pkg/app"
)

type Api struct {
	App *app.App
}

func New(app *app.App) *Api {
	return &Api{
		App: app,
	}
}

func (api *Api) SetupRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/baskets", func(r chi.Router) {
		r.Post("/", api.CreateBasket)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(api.basketCtx)

			r.Put("/", api.CloseBasket)
			r.Put("/items", api.AddItem)
		})
	})

	router.Route("/items", func(r chi.Router) {
		r.Get("/", api.GetItems)
	})

	return router
}

func (api *Api) basketCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		basket, err := api.App.DB.BasketsStore.Get(id)
		if err != nil {
			switch err.(type) {
			case *db.RecordNotFoundError:
				handleError(&BasketNotFoundError{
					ID: id,
				}, w)
			default:
				handleError(err, w)
			}

			return
		}

		ctx := context.WithValue(r.Context(), "basket", basket)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleError(err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}

	log.Printf("Handling error %+v", err)

	switch err.(type) {
	case *BasketNotFoundError:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	bs, _ := json.Marshal(map[string]interface{}{"error": err.Error()})
	w.Write(bs)

	return true
}

func writeJSON(w http.ResponseWriter, data []byte, status int) {
	w.WriteHeader(status)
	w.Header().Add("content-type", "application/json")

	w.Write(data)
}
