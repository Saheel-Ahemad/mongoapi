// 1. List down the requirements such as local mongo uri, database name and collection name
// 2. Initialise the collection from mongo by reference
// 3. With init function, set client options for created uri
// 4. Check for mongo connection with the created client options
// 5. Log error if any
// 6. Access the collection by passing database name and collection name

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection string
const connectionString = "mongodb://localhost:27017"

const dbName = "netflix"
const collectionName = "watchlist"

// Global variable for collection reference
var collection *mongo.Collection

// Initialize MongoDB connection
func init() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set MongoDB client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping MongoDB to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	// Get collection reference
	collection = client.Database(dbName).Collection(collectionName)

	log.Println("Successfully connected to MongoDB!")
}

// package controller

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gorilla/mux"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// const connectionString = "mongodb://localhost:27017"
// const dbName = "netflix"
// const collectionName = "watchlist"

// // Most important for connection
// var collection *mongo.Collection

// // Connect with mongodb
// func init() {
// 	// Create a context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Set client options
// 	clientOptions := options.Client().ApplyURI(connectionString)

// 	// Connect to MongoDB
// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Fatal("Failed to connect to MongoDB:", err)
// 	}

// 	// Ping the database to verify connection
// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		log.Fatal("Failed to ping MongoDB:", err)
// 	}

// 	// Get collection reference
// 	collection = client.Database(dbName).Collection(collectionName)

// 	log.Println("Successfully connected to MongoDB!")
// }

// Mongodb helpers

func insertOneMovie(movie string) {
	insertedMovie, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted 1 movie to the database with id", insertedMovie.InsertedID)
}

// Updation of the record
func updateOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"Watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count", result.ModifiedCount)
}

// Deletion one of the record
func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal()
	}
	fmt.Println("Deleted record", deleteCount)
}

// Deletion all records from mongodb
func deleteAllMovie() int64 {
	deleteCount, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal()
	}
	fmt.Println("The number of movies deleted", deleteCount.DeletedCount)
	return deleteCount.DeletedCount
}

// Get all movies from mongodb
func getAllMovies() []primitive.M {
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal()
	}
	var movies []primitive.M
	for cursor.Next(context.Background()) {
		var movie bson.M
		err := cursor.Decode(&movie)
		if err != nil {
			log.Fatal()
		}
		movies = append(movies, movie)
	}
	defer cursor.Close(context.Background())
	return movies
}

func GetMyAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoder")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoder")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie string
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoder")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoder")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)
}
