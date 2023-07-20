package handlers

import (
	"fmt"
	"log"
	"net/http"
	"project/internal/web/templates"
)

type cache interface {
	GetOrder(orederUID string) ([]byte, error)
}

func HandleMain(cache cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := templates.GetMainTemplate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		page := templates.MainPage{}
		orderID := r.FormValue("orderID")
		if orderID == "" {
			tmpl.Execute(w, "Order ID is empty")
			return
		}

		page.OrderJSON, err = cache.GetOrder(orderID)
		if err != nil {
			log.Println(err)
			tmpl.Execute(w, fmt.Sprintf("Order with ID %s not found", orderID))
			return
		}

		tmpl.Execute(w, page.OrderJSON)
	}
}
