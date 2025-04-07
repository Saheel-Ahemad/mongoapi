// 1. List down the requirements such as local mongo uri, database name and collection name
// 2. Initialise the collection from mongo by reference
// 3. With init function, set client options for created uri
// 4. Check for mongo connection with the created client options
// 5. Log error if any
// 6. Access the collection by passing database name and collection name

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mongoapi/model"
	"net/http"
	"time"

	// "github.com/gorilla/mux"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection string
// const connectionString = "mongodb://localhost:27017"
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
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB`
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

// Home function to handle the root route
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<h1>Welcome to my Netflix API</h1>")
}

// Movie struct to represent a movie document in MongoDB
type Movie struct {
	Title    string
	Watched  bool
	Released int
}

// Function to insert a movie into the database
func InsertInterstellar(w http.ResponseWriter, r *http.Request) {
	movie := model.Netflix{
		ID:      primitive.NewObjectID(),
		Movie:   "Interstellar",
		Watched: false,
	}
	insertedMovie, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, "<h1>Something has went wrong</h1>")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Interstellar Movie inserted successfully!",
		"id":      insertedMovie.InsertedID,
	})

	// w.Header().Set("Content-Type", "text/html")
	// fmt.Fprintln(w, "<h1>Interstellar</h1>")
}

// // Updation of the record
// func updateOneMovie(movieId string) {
// 	id, _ := primitive.ObjectIDFromHex(movieId)
// 	filter := bson.M{"_id": id}
// 	update := bson.M{"$set": bson.M{"Watched": true}}

// 	result, err := collection.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Modified count", result.ModifiedCount)
// }

// // Deletion one of the record
// func deleteOneMovie(movieId string) {
// 	id, _ := primitive.ObjectIDFromHex(movieId)
// 	filter := bson.M{"id": id}
// 	deleteCount, err := collection.DeleteOne(context.Background(), filter)

// 	if err != nil {
// 		log.Fatal()
// 	}
// 	fmt.Println("Deleted record", deleteCount)
// }

// // Deletion all records from mongodb
// func deleteAllMovie() int64 {
// 	deleteCount, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
// 	if err != nil {
// 		log.Fatal()
// 	}
// 	fmt.Println("The number of movies deleted", deleteCount.DeletedCount)
// 	return deleteCount.DeletedCount
// }

// Primitive datatype stores the slice of maps
//	{
//	  {
//		"_id": "ObjectID",
//		"title": "Interstellar",
//		"watched": true,
//		"released": 2014
//	  }
//	  {
//		"_id": "ObjectID2",
//		"title": "Inception",
//		"watched": false,
//		"released": 2010
//	  }
//	}

// Get all movies from mongodb
func GetAllMovies() []primitive.M {
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

// Retrive all movies data
// Usage of x-www-form-urlencoder due to accepting the data in the form of key value pair

func GetMyAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoder")
	allMovies := GetAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

// // Create a single movie
// func CreateMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/x-www-form-urlencoder")
// 	w.Header().Set("Allow-Control-Allow-Methods", "POST")

// 	var movie Movie
// 	_ = json.NewDecoder(r.Body).Decode(&movie)
// 	// insertOneMovie(movie)
// 	json.NewEncoder(w).Encode(movie)
// }

// // Mark a movie as watched
// func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
// 	// var movie Movie
// 	// insertOneMovie(movie)
// 	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

// 	params := mux.Vars(r)
// 	updateOneMovie(params["id"])
// 	json.NewEncoder(w).Encode(params["id"])
// }

// // Delete a movie from record
// func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
// 	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

// 	params := mux.Vars(r)
// 	deleteOneMovie(params["id"])
// 	json.NewEncoder(w).Encode(params["id"])
// }

// // Delete all movies from record
// func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/x-www-form-urlencoder")
// 	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

// 	count := deleteAllMovie()
// 	json.NewEncoder(w).Encode(count)
// }
