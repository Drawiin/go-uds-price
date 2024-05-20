package main

import (
	"net/http"
	"github.com/drawiin/go-usd-price/server/price"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /price", price.PriceController{})
	http.ListenAndServe(":8080", mux)
}
