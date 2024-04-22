package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/velvetreactor/metrics/pkg/ynabcollector"
)

func main() {
	hourlyRegistry := prometheus.NewRegistry()

	ynabCollector := &ynabcollector.YNABCollector{}
	hourlyRegistry.MustRegister(ynabCollector)

	r := chi.NewRouter()

	r.Route("/transactions", func(r chi.Router) {
	})

	r.Route("/notes", func(r chi.Router) {
		r.Get("/new", func(w http.ResponseWriter, r *http.Request) {
			file, err := os.ReadFile("views/add_note.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			data := struct {
				Title    string
				Header   string
				Contents template.HTML
			}{
				Title:    "Add note",
				Header:   "Add note",
				Contents: template.HTML(file),
			}

			tmpl, err := template.ParseFiles("views/layout.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	})

	// Metrics
	r.Get(
		"/ynab/metrics",
		promhttp.HandlerFor(
			hourlyRegistry,
			promhttp.HandlerOpts{},
		).(http.HandlerFunc),
	)

	port := os.Getenv("METRICS_SERVER_PORT")
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatal(err)
	}
}
