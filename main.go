package main

import (
	"crud/controller"
	"crud/database"
	"crud/repository"
	"crud/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	personRepository repository.PersonRepository = repository.NewPersonRepository(database.MongoDB())
	personService    service.PersonService       = service.NewPersonService(personRepository)
	personController controller.PersonController = controller.NewPersonController(personService)
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/people", FindAll).Methods("GET")
	router.HandleFunc("/people/{id}", FindById).Methods("GET")
	router.HandleFunc("/people", Save).Methods("POST")
	router.HandleFunc("/people/{id}", Update).Methods("PUT")
	router.HandleFunc("/people/{id}", Remove).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":80", router))
}

func FindAll(w http.ResponseWriter, r *http.Request) {
	people := personController.FindAll()
	w.Write(toJson(people))
}

func FindById(w http.ResponseWriter, r *http.Request) {
	person := personController.FindById(w, r)
	w.Write(toJson(person))
}

func Save(w http.ResponseWriter, r *http.Request) {
	id := personController.Save(w, r)
	w.Write(toJson(id))
}

func Update(w http.ResponseWriter, r *http.Request) {
	id := personController.Update(w, r)
	w.Write(toJson(id))
}

func Remove(w http.ResponseWriter, r *http.Request) {
	id := personController.Remove(w, r)
	w.Write(toJson(id))
}

func toJson(person interface{}) []byte {
	j, err := json.Marshal(person)
	if err != nil {
		log.Fatal(err.Error())
	}
	return j
}