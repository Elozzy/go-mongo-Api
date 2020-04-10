package main

import (
	"fmt"
	"context"
	"net/http"
	"time"
	"encoding/json"
	"log"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gorilla/mux"
)



//DataBase Name = Village
var client *mongo.Client

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}




//delete person
// func DeletePersonEndpoint(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	id, _ := primitive.ObjectIDFromHex(params["id"])
// 	var person Person
// 	collection := client.Database("villagepeople").Collection("people")
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	for index, item := range Person {
// 		if item.ID == params["id"]{
// 			person = append(Person[:index], Person[index+1:]...)
// 			break
// 		}
// 	}
// 	json.NewEncoder(w).Encode(person)
	

// }


func main() {
	fmt.Println("Starting the application....#######")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()

	//Endpoints routes
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")
	//setting port addres
	log.Fatal(http.ListenAndServe(":5000", router))

}

//Get a singular person on the data base

func GetPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection := client.Database("villagepeople").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(person)
}

//GetRequest/ all user on the db
func GetPeopleEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var people []Person
	collection := client.Database("villagepeople").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(people)
}

//post/created a user
func CreatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(r.Body).Decode(&person)
	collection := client.Database("villagepeople").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(w).Encode(result)

}