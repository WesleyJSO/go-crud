package controller

import (
	"crud/entity"
	"crud/service"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type PersonController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	FindById(w http.ResponseWriter, r *http.Request)
	Remove(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	service service.PersonService
}

func NewPersonController(service service.PersonService) PersonController {
	return &controller{service}
}

func (c *controller) FindAll(w http.ResponseWriter, r *http.Request) {
	people := c.service.FindAll()
	w.Write(toJson(people))
}

func (c *controller) FindById(w http.ResponseWriter, r *http.Request) {
	person := c.service.FindById(mux.Vars(r)["id"])
	if (person == entity.Person{}) {
		http.NotFound(w, r)
	}
	w.Write(toJson(person))
}

func (c *controller) Remove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	deleteCount := c.service.Remove(id)
	if deleteCount == 0 {
		log.Print("No entity found for id: " + id)
	}
	w.Write(toJson(id))
}

func (c *controller) Save(w http.ResponseWriter, r *http.Request) {
	var body []byte = readBodyFromRequest(r)
	var person entity.Person
	json.Unmarshal(body, &person)
	id := c.service.Save(&person)
	w.Write(toJson(id))

}

func (c *controller) Update(w http.ResponseWriter, r *http.Request) {
	var body []byte = readBodyFromRequest(r)
	id := c.service.Update(mux.Vars(r)["id"], body)
	w.Write(toJson(id))
}

func readBodyFromRequest(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Read body error: ", err.Error())
	}
	return body
}

func toJson(person interface{}) []byte {
	j, err := json.Marshal(person)
	if err != nil {
		log.Fatal(err.Error())
	}
	return j
}
