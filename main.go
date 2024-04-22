package main

import (
	"fmt"
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
