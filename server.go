package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/streadway/amqp"
)

type Router interface {
	Get(pattern string, handlerFunc http.HandlerFunc)
}

type Server struct {
	router    Router
	service   StockService
	upstreams []Endpoint
}

func NewServer(r Router, brokerConn *amqp.Channel) *Server {
	return &Server{
		router:  r,
		service: NewStockHandler(brokerConn),
		upstreams: []Endpoint{
			GET_CSV_FOR_STOCK,
		},
	}
}

func (s *Server) ServeHTTP() {
	s.router.Get("/{stock}", func(w http.ResponseWriter, r *http.Request) {

		p := chi.URLParam(r, "stock")
		queryParams := r.URL.Query()

		sender := queryParams.Get("sender")
		if sender == "" {
			http.Error(w, "query param sender is mandatory", http.StatusBadRequest)
			return
		}
		if p == "" {
			http.Error(w, "invalid url parameter provided", http.StatusBadRequest)
			return
		}
		fmt.Printf("%s&s=%s", string(GET_CSV_FOR_STOCK), p)
		resp, err := http.Get(fmt.Sprintf("%s&s=%s", string(GET_CSV_FOR_STOCK), p))
		if err != nil {
			fmt.Println(err)
			return
		}

		s.service.ParseStock(resp.Body, sender)

	})
}
