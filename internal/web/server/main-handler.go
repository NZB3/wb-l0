package server

import (
	"fmt"
	"log"
	"net/http"
	"project/internal/web/templates"
)

func (s *server) handleMain() http.HandlerFunc {
	const op = "web.server.handleMain"

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := templates.GetMainTemplate()
		if err != nil {
			log.Printf("%s: %s", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if r.Method != http.MethodPost {
			log.Printf("%s: Method %s", op, r.Method)
			tmpl.Execute(w, nil)
			return
		}

		page := templates.MainPage{}
		orderID := r.FormValue("orderID")
		if orderID == "" {
			log.Printf("%s: Order ID is empty", op)
			tmpl.Execute(w, "Order ID is empty")
			return
		}
		log.Println(orderID)

		page.OrderJSON, err = s.cache.GetOrder(orderID)
		if err != nil {
			log.Printf("%s: %s", op, err)
			tmpl.Execute(w, fmt.Sprintf("Order with ID %s not found", orderID))
			return
		}
		log.Printf("%s: %s", op, page.OrderJSON)

		tmpl.Execute(w, string(page.OrderJSON))
	}
}
