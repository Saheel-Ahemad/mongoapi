package router

import (
	"mongoapi/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	// log.Fatal(http.ListenAndServe(":8080", nil))
	router := mux.NewRouter()
	router.HandleFunc("/", controller.Home).Methods("GET")
	router.HandleFunc("/add-interstellar", controller.InsertInterstellar).Methods("GET")
	router.HandleFunc("/all-movies", controller.GetMyAllMovies).Methods("GET")
	// router.HandleFunc("/api/movie", controller.CreateMovie).Methods("POST")
	// router.HandleFunc("/api/movie/{id}", controller.MarkAsWatched).Methods("PUT")
	// router.HandleFunc("/api/movie/{id}", controller.DeleteAMovie).Methods("DELETE")
	// router.HandleFunc("/api/deleteallmovie", controller.DeleteAllMovies).Methods("PUT")

	return router
}
