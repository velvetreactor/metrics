package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/velvetreactor/metrics/pkg/ynabcollector"
)

func main() {
	dbConn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Printf("unable to connect to database: %v", err)
		os.Exit(1)
	}
	defer dbConn.Close(context.Background())

	hourlyRegistry := prometheus.NewRegistry()

	ynabCollector := &ynabcollector.YNABCollector{}
	hourlyRegistry.MustRegister(ynabCollector)

	r := chi.NewRouter()

	r.Route("/transactions", func(r chi.Router) {
	})

	r.Route("/notes", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			commandTag, err := dbConn.Exec(
				context.Background(),
				"INSERT INTO notes (body) VALUES ($1)",
				r.Form.Get("body"),
			)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if commandTag.RowsAffected() != 1 {
				log.Printf("Expected to affect 1 row, affected: %d", commandTag.RowsAffected())
				return
			}
		})

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
