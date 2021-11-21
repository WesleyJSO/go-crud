package service

import (
	"crud/entity"
	"crud/repository"
	"encoding/json"
)

type PersonService interface {
	FindAll() []entity.Person
	FindById(id string) entity.Person
	Remove(id string) int64
	Save(person *entity.Person) string
	Update(id string, body []byte) string
}

type service struct {
	repository repository.PersonRepository
}

func NewPersonService(repository repository.PersonRepository) PersonService {
	return &service{repository}
}

func (s *service) FindAll() []entity.Person {
	return s.repository.FindAll()
}

func (s *service) FindById(id string) entity.Person {
	return s.repository.FindById(id)
}

func (s *service) Remove(id string) int64 {
	return s.repository.Remove(id)
}

func (s *service) Save(person *entity.Person) string {
	return s.repository.Save(person)
}

func (s *service) Update(id string, body []byte) string {
	var person entity.Person = s.FindById(id)
	json.Unmarshal(body, &person)
	return s.repository.Update(id, &person)
}
