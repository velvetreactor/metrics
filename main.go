package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	LineItemAssignedDesc = prometheus.NewDesc(
		"line_item_assigned",
		"The monthly amount assigned to this line item",
		[]string{"name", "category_group_name"},
		nil,
	)

	LineItemUsedDesc = prometheus.NewDesc(
		"line_item_used",
		"The amount used for this line item",
		[]string{"name", "category_group_name"},
		nil,
	)
)

type YNABCollector struct{}

func (c *YNABCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- LineItemAssignedDesc
	ch <- LineItemUsedDesc
}

func (c *YNABCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("Collecting YNAB metrics...")

	budgetId := os.Getenv("MY_BUDGET_ID")
	if budgetId == "" {
		log.Println("MY_BUDGET_ID is blank")
		return
	}

	ynabToken := os.Getenv("YNAB_TOKEN")
	if ynabToken == "" {
		log.Println("YNAB_TOKEN is empty")
		return
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("https://api.ynab.com/v1/budgets/%s/categories", budgetId),
		nil,
	)

	if err != nil {
		log.Printf("error creating request: %v", err)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ynabToken))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var apiResp Response
	err = json.Unmarshal(b, &apiResp)
	if err != nil {
		log.Printf("error unmarshalling json response: %v", err)
	}

	for _, group := range apiResp.Data.CategoryGroups {
		if group.Hidden {
			continue
		}

		for _, lineItem := range group.Categories {
			if lineItem.Hidden {
				continue
			}

			ch <- prometheus.MustNewConstMetric(
				LineItemAssignedDesc,
				prometheus.GaugeValue,
				lineItem.GetBudgeted(),
				lineItem.Name,
				lineItem.GroupName,
			)

			ch <- prometheus.MustNewConstMetric(
				LineItemUsedDesc,
				prometheus.GaugeValue,
				lineItem.GetUsed(),
				lineItem.Name,
				lineItem.GroupName,
			)
		}
	}
}

type Response struct {
	Data Data
}

type CategoryGroup struct {
	Hidden     bool      `json:"hidden"`
	Categories LineItems `json:"categories"`
}

type Data struct {
	CategoryGroups []CategoryGroup `json:"category_groups"`
}

type LineItems []LineItem
type LineItem struct {
	// Category
	Name      string `json:"name"`
	Budgeted  int    `json:"budgeted"`
	Hidden    bool   `json:"hidden"`
	GroupName string `json:"category_group_name"`
	Activity  int    `json:"activity"`
}

func (l LineItem) GetBudgeted() float64 {
	return float64(l.Budgeted) / 1000
}

func (l LineItem) GetUsed() float64 {
	return -(float64(l.Activity) / 1000)
}

func main() {
	hourlyRegistry := prometheus.NewRegistry()

	ynabCollector := &YNABCollector{}
	hourlyRegistry.MustRegister(ynabCollector)

	http.Handle("/ynab/metrics", promhttp.HandlerFor(hourlyRegistry, promhttp.HandlerOpts{}))

	fmt.Println("Server starting...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
