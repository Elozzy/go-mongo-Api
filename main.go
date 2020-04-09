package main

import (
	"fmt"
	"context"
	"net/http"
	"time"
	"encoding/json"
	// "log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"
)


// Define models
type Person struct{
	ID 			primitive.ObjectID	    `json:"_id,"pmitempty" bson:"_id,omitempty"`
	Firstname 	string 					`json:"firstname, omitempty" "bsona;"firstname, omitempty"`
	Lastname 	string 					`json:"lastname, omitempty" "bsona;"lastname, omitempty"`
}


var client *mongo.Client



//post

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("villagepeople").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)

}


func main() {
	fmt.Println("Starting the application")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, "mongodb://localhost:27017")
	router := mux.NewRouter()

	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")


	http.ListenAndServe(":5000", router)

}