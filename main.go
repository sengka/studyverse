package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	// Anasayfa
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/greeting.html")
	}).Methods("GET")

	log.Println("Server çalışıyor : http://localhost:6060")
	log.Fatal(http.ListenAndServe(":6060", r))
}
