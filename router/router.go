package router

import (
	"github.com/dedeTr/url-shortener/middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/newurl", middleware.CreateURL).Methods("POST", "OPTIONS")
	router.HandleFunc("/search/{url_req}", middleware.Geturl).Methods("GET", "OPTIONS")

	return router
}
