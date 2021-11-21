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
	FindAll() []entity.Person
	FindById(w http.ResponseWriter, r *http.Request) entity.Person
	Remove(w http.ResponseWriter, r *http.Request) string
	Save(w http.ResponseWriter, r *http.Request) string
	Update(w http.ResponseWriter, r *http.Request) string
}

type controller struct {
	service service.PersonService
}

func NewPersonController(service service.PersonService) PersonController {
	return &controller{service}
}

func (c *controller) FindAll() []entity.Person {
	return c.service.FindAll()
}

func (c *controller) FindById(w http.ResponseWriter, r *http.Request) entity.Person {
	person := c.service.FindById(mux.Vars(r)["id"])
	if (person == entity.Person{}) {
		http.NotFound(w, r)
	}
	return person
}

func (c *controller) Remove(w http.ResponseWriter, r *http.Request) string {
	id := mux.Vars(r)["id"]
	deleteCount := c.service.Remove(id)
	if deleteCount == 0 {
		log.Print("No entity found for id: " + id)
	}
	return id
}

func (c *controller) Save(w http.ResponseWriter, r *http.Request) string {
	var body []byte = readBodyFromRequest(r)
	var person entity.Person
	json.Unmarshal(body, &person)
	return c.service.Save(&person)
}

func (c *controller) Update(w http.ResponseWriter, r *http.Request) string {
	var body []byte = readBodyFromRequest(r)
	return c.service.Update(mux.Vars(r)["id"], body)
}

func readBodyFromRequest(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Read body error: ", err.Error())
	}
	return body
}
