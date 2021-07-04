package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"message": "Hello Docker!",
		}
		json.NewEncoder(rw).Encode(response)
	})

	log.Println("Server is running!")
	http.ListenAndServe(":4000", router)
}
