package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Saheel-Ahemad/mongoapi/router"
)

func main() {
	fmt.Println("MongoDB API")
	r := router.Router()
	fmt.Println("Server is getting started...")
	log.Fatal(http.ListenAndServe(":27017", r))
	fmt.Println("Listening at port 27017 ...")
}