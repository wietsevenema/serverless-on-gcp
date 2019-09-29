package main

import (
	"chapter-4/shared/config"
	"chapter-4/shared/database"
	"chapter-4/shared/environment"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type env struct {
	*environment.Env
}

func main() {
	// Initialize services
	env := &env{environment.NewEnv()}
	err := env.InitDB()
	if err != nil {
		panic(err)
	}
	err = env.InitTasks()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", env.serveIndex)
	http.HandleFunc("/booking", env.postBooking)

	// The trailing slash means the entire subtree is matched
	// Example: /booking/123 will go to env.getBooking
	http.HandleFunc("/booking/", env.getBooking)

	// Start server
	log.Println("Listening on port " + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))

}

// serveIndex returns the index.html file
func (e *env) serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

// postBooking receives a form POST, save the fields in the database
// and send a tasks to the booking function
func (e *env) postBooking(w http.ResponseWriter, r *http.Request) {
	// Get the form values from the request
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error parsing form",
			http.StatusInternalServerError)
	}
	// Check if the from is not empty
	if len(r.PostForm) > 0 {
		// Initialize the booking object
		booking := &database.Booking{
			Name:  r.PostForm.Get("Name"),
			Email: r.PostForm.Get("Email"),
		}

		err := e.submitBooking(r.Context(), booking)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error saving booking",
				http.StatusInternalServerError)
		}

		// Redirect to the booking detail page
		http.Redirect(
			w, r,
			fmt.Sprintf("/booking/%s", booking.ID),
			http.StatusFound,
		)
	}
}

func (e *env) submitBooking(ctx context.Context, booking *database.Booking) error {
	// Save the booking object
	err := e.Db.SaveBooking(ctx, booking)

	if err != nil {
		return fmt.Errorf("error saving booking: %v", err)
	}

	// Send a task to the booking function
	err = e.Tasks.AddTask(ctx,
		"bookings",
		"background",
		booking.ID,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error submitting task: %v", err)
	}
	return nil
}

// getBooking renders an HTML template with the booking info.
func (e *env) getBooking(w http.ResponseWriter, r *http.Request) {
	// First, fetch the ID from the path.
	// Example: /booking/123 => ID: 123
	ID := strings.TrimPrefix(r.URL.Path, "/booking/")

	// Minor sanitization of the id, you want to be more
	// strict in a production setting
	ID = strings.ReplaceAll(ID, "/", "")

	if ID == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Fetch the booking from the database by ID
	result, err := e.Db.GetBooking(r.Context(), ID)
	if err != nil {
		fail("error getting booking", err, w)
		return
	}

	if !result.Exists() {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Render booking in template
	tpl := template.Must(
		template.ParseFiles("web/booking.html"))
	tpl.Execute(w, result)

}

// fail is a helper that prints a message to the page and
// logs the full error message.
func fail(message string, err error, w http.ResponseWriter) {
	log.Println(message + ": " + err.Error())
	http.Error(w,
		message,
		http.StatusInternalServerError)
}
