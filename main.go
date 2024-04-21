package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/velvetreactor/metrics/pkg/ynabcollector"
)

func main() {
	hourlyRegistry := prometheus.NewRegistry()

	ynabCollector := &ynabcollector.YNABCollector{}
	hourlyRegistry.MustRegister(ynabCollector)

	// Metrics
	http.Handle("/ynab/metrics", promhttp.HandlerFor(hourlyRegistry, promhttp.HandlerOpts{}))

	// Pages
	http.HandleFunc("/transactions/new", func(w http.ResponseWriter, req *http.Request) {
		data, err := os.ReadFile("views/add_transaction.html")
		if err != nil {
			log.Printf("Error reading file: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		w.Write(data)
	})

	port := os.Getenv("METRICS_SERVER_PORT")
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}
