package repository

import (
	"context"
	"go-crud/database"
	"go-crud/entity"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const COLL_PEOPLE = "people"

type PersonRepository interface {
	FindAll() []entity.Person
	FindById(id string) entity.Person
	Remove(id string) int64
	Save(person *entity.Person) string
	Update(id string, person *entity.Person) string
}

type repository struct {
	database database.Database
}

func NewPersonRepository(database database.Database) PersonRepository {
	return &repository{database}
}

func (r *repository) FindAll() []entity.Person {
	cur, err := r.database.GetCollection(COLL_PEOPLE).Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal("Fetch people error: ", err.Error())
	}

	people := make([]entity.Person, 0)
	for cur.Next(context.Background()) {
		var person entity.Person
		cur.Decode(&person)
		people = append(people, person)
	}
	return people
}

func (r *repository) FindById(id string) entity.Person {
	var person entity.Person
	r.database.GetCollection(COLL_PEOPLE).FindOne(context.Background(), bson.M{"_id": hexId(id)}).Decode(&person)
	return person
}

func (r *repository) Remove(id string) int64 {
	res, err := r.database.GetCollection(COLL_PEOPLE).DeleteOne(context.Background(), bson.M{"_id": hexId(id)})
	if err != nil {
		log.Fatal("Delete error: ", err.Error())
	}
	return res.DeletedCount
}

func (r *repository) Save(person *entity.Person) string {
	_, err := r.database.GetCollection(COLL_PEOPLE).InsertOne(context.Background(), person)
	if err != nil {
		log.Fatal("Insert error: ", err.Error())
	}
	return person.Name
}

func (r *repository) Update(id string, person *entity.Person) string {
	update := bson.M{"$set": bson.M{"name": person.Name, "age": person.Age}}
	_, err := r.database.GetCollection(COLL_PEOPLE).UpdateByID(context.Background(), hexId(id), update)
	if err != nil {
		log.Fatal("Fetch by id error: ", err.Error())
	}
	return person.Id.String()
}

func hexId(pathParam string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(pathParam)
	if err != nil {
		log.Fatal("Invalid id error: ", err.Error())
	}
	return id
}
