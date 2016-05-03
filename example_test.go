package serveseq_test

import (
	"fmt"
	"github.com/fatlotus/serveseq"
	"net/http"
)

// If the user is not logged in, redirect to the login page.
func RequireLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/login" {
		fmt.Fprintf(w, "[login form]")

	} else if r.FormValue("session") == "" {
		http.Redirect(w, r, "/login", 307)
	}
}

// Capture any unhandled requests.
func CustomNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "not found :(")
}

func Example() {
	// Create a conventional application.
	mux := http.NewServeMux()
	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world!")
	})

	// Compose it with some middleware.
	seq := serveseq.New()
	seq.NextFunc(RequireLogin)
	seq.Next(mux)
	seq.NextFunc(CustomNotFound)

	RunRequests(seq, "/app", "/app?session=yes", "/notapp?session=yes")

	// Output:
	// /app: [login form]
	// /app?session=yes: hello world!
	// /notapp?session=yes: not found :(
}
