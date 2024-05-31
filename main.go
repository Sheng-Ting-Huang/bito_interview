package main

import (
	"log"
	"net/http"

	"github.com/bito_interview/api"
)

func main() {
	log.Println("listen to :8080")
	log.Fatal(http.ListenAndServe(":8080", api.NewRouter()))
}
