package server

import (
	"fmt"
	"net/http"
)

type cache interface {
	GetOrder(orederUID string) ([]byte, error)
}

type server struct {
	cache cache
}

func NewServer(cache cache) *server {
	return &server{cache: cache}
}

func (s *server) RunServer(cache cache) {
	fmt.Println("Web Server running at http://localhost:8080")
	http.HandleFunc("/", s.handleMain())
	http.ListenAndServe(":8080", nil)
}
