package router

import (
	"crud/controller"
	"crud/database"
	"crud/repository"
	"crud/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var db database.Database = database.MongoDB()

func personRoutes(router *mux.Router) {

	personRepository := repository.NewPersonRepository(db)
	personService := service.NewPersonService(personRepository)
	personController := controller.NewPersonController(personService)

	router.HandleFunc("/people", personController.FindAll).Methods("GET")
	router.HandleFunc("/people/{id}", personController.FindById).Methods("GET")
	router.HandleFunc("/people", personController.Save).Methods("POST")
	router.HandleFunc("/people/{id}", personController.Update).Methods("PUT")
	router.HandleFunc("/people/{id}", personController.Remove).Methods("DELETE")
}

func Init() {
	router := mux.NewRouter()

	personRoutes(router)

	log.Fatal(http.ListenAndServe(":80", router))
}
