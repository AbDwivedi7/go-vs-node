package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Crud struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Published   bool               `json:"published,omitempty" bson:"published,omitempty"`
}

var client *mongo.Client

func postCrud(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var crud Crud
	json.NewDecoder(r.Body).Decode(&crud)
	fmt.Println(crud)
	collection := client.Database("test").Collection("cruds")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, crud)
	json.NewEncoder(w).Encode(result)
}

func getCrud(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var cruds []Crud
	collection := client.Database("test").Collection("cruds")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `}"`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var crud Crud
		cursor.Decode(&crud)
		cruds = append(cruds, crud)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `}"`))
		return
	}
	json.NewEncoder(w).Encode(cruds)
}

func getACrud(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var crud Crud
	collection := client.Database("test").Collection("cruds")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, Crud{ID: id}).Decode(&crud)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `}"`))
		return
	}
	json.NewEncoder(w).Encode(crud)
}

func main() {
	fmt.Println("Starting the application....")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	password := os.Getenv("PASSWORD")
	connectionString := "mongodb+srv://abdwivedi:" + password + "@testing.si9qg.mongodb.net/test?retryWrites=true&w=majority"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	router := mux.NewRouter()
	router.HandleFunc("/api/crud/postCrud", postCrud).Methods("POST")
	router.HandleFunc("/api/crud/getCrud", getCrud).Methods("GET")
	router.HandleFunc("/api/crud/getCrud/{id}", getACrud).Methods("GET")
	router.HandleFunc("/api/crud/updateCrud", postCrud).Methods("PUT")
	router.HandleFunc("/api/crud/deleteCrud", postCrud).Methods("DELETE")

	port := os.Getenv("PORT")
	http.ListenAndServe(port, router)

}
