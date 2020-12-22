package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func testFirst(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Response")
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/test", testFirst)
	log.Fatal(http.ListenAndServe(":8080", router))
}
