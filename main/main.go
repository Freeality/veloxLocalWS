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
	contaController.CriarContasParaTeste()

	router.HandleFunc("/conta/delete/{id}", contaController.requestDelete).Methods("DELETE")
	router.HandleFunc("/conta/{id}", contaController.requestGet).Methods("GET")
	router.HandleFunc("/conta/post/", contaController.requestPost).Methods("POST")

	return router
}
