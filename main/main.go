package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"

	"github.com/freeality/veloxLocalWS/controllers"
)

func main() {

	log.Fatal(http.ListenAndServe(":8080", registraPaths()))
}

func registraPaths() http.Handler {

	router := mux.NewRouter().StrictSlash(true)
	contaController := controllers.NewContaController()

	router.HandleFunc("/conta/delete/{id}", contaController.WebDelete).Methods("DELETE")

	return router
}
