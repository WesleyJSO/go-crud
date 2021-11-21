package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const COLL_PEOPLE = "people"

var client *mongo.Client

func main() {
	client = Connect(GetEnv("DBCONN"))

	router := mux.NewRouter()
	router.HandleFunc("/", FindAll).Methods("GET")
	router.HandleFunc("/{id}", FindById).Methods("GET")
	router.HandleFunc("/", Save).Methods("POST")
	router.HandleFunc("/{id}", Update).Methods("PUT")
	router.HandleFunc("/{id}", Remove).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":80", router))
}

// STRUCT
type Person struct {
	Id   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty"`
	Age  int                `json:"age,omitempty"`
}

// STRUCT

//HANDLERS
func FindAll(w http.ResponseWriter, r *http.Request) {
	cur, err := GetCollection(COLL_PEOPLE).Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal("Fetch people error: ", err.Error())
	}

	people := make([]Person, 0)
	for cur.Next(context.Background()) {
		var person Person
		cur.Decode(&person)
		people = append(people, person)
	}
	w.Write(ToJson(people))
}

func FindById(w http.ResponseWriter, r *http.Request) {
	id := HexId(mux.Vars(r)["id"])
	var person Person
	GetCollection(COLL_PEOPLE).FindOne(context.Background(), bson.M{"_id": id}).Decode(&person)
	if (Person{} == person) {
		http.NotFound(w, r)
	}
	w.Write(ToJson(person))
}

func Save(w http.ResponseWriter, r *http.Request) {
	var body []byte = ReadBodyFromRequest(r)
	var person Person
	json.Unmarshal(body, &person)

	res, err := GetCollection(COLL_PEOPLE).InsertOne(context.Background(), person)
	if err != nil {
		log.Fatal("Insert error: ", err.Error())
	}
	w.Write(ToJson(res))
}

func Update(w http.ResponseWriter, r *http.Request) {
	var body []byte = ReadBodyFromRequest(r)
	id := HexId(mux.Vars(r)["id"])
	var person Person

	GetCollection(COLL_PEOPLE).FindOne(context.Background(), bson.M{"_id": id}).Decode(&person)
	json.Unmarshal(body, &person)
	
	update := bson.M{"$set": bson.M{"name": person.Name, "age": person.Age}}

	_, err := GetCollection(COLL_PEOPLE).UpdateByID(context.Background(), id, update)
	if err != nil {
		log.Fatal("Fetch by id error: ", err.Error())
	}
	w.Write(ToJson(person))
}

func Remove(w http.ResponseWriter, r *http.Request) {
	id := HexId(mux.Vars(r)["id"])
	res, err := GetCollection(COLL_PEOPLE).DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Fatal("Delete error: ", err.Error())
	}
	w.Write(ToJson(res))
}
//HANDLERS

// MONGODB
func GetCollection(name string) *mongo.Collection {
	return client.Database(GetEnv("DBNAME")).Collection(name)
}

func Connect(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Mongo client read error: ", err.Error())
	}

	fmt.Println(client.Database(GetEnv("DBNAME")).ListCollectionNames(
		ctx,
		bson.M{},
	))
	return client
}

// MONGODB

// ENVIRONMENT
func GetEnv(value string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(value)
}

// ENVIRONMENT

// UTILS
func HexId(pathParam string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(pathParam)
	if err != nil {
		log.Fatal("Invalid id error: ", err.Error())
	}
	return id
}

func ReadBodyFromRequest(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Read body error: ", err.Error())
	}
	return body
}

func ToInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err.Error())
	}
	return val
}

func ToJson(person interface{}) []byte {
	j, err := json.Marshal(person)
	if err != nil {
		log.Fatal(err.Error())
	}
	return j
}

// UTILS
