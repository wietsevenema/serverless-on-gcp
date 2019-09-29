package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"chapter-4/shared/config"
	"chapter-4/shared/environment"
)

type env struct {
	*environment.Env
}

func main() {
	// Grab the PORT from the environment
	log.Println("Listening on port " + config.Port)

	env := &env{environment.NewEnv()}
	err := env.InitDB()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", env.processBooking)

	// Start the HTTP server on PORT
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

func (e *env) processBooking(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ID := strings.ReplaceAll(r.URL.Path, "/", "")

	if ID == "" {
		fmt.Fprintln(w, "No ID present")
		return
	}

	booking, err := e.Db.GetBooking(r.Context(), ID)
	if err != nil {
		log.Printf("Error handling [%v]: %+v", ID, err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if booking.Exists() {
		// This is where we apply a few processing steps
		// in the real-world
		// ...

		// Set Processed to true and save the result back
		// to the database
		booking.Processed = true
		e.Db.SaveBooking(r.Context(), booking)
		return
	}
	http.Error(w, "Not found", http.StatusNotFound)

}
