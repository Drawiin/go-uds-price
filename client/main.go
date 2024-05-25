package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type PriceResponse struct {
	Price float64 `json:"price"`
}

func main() {
	price, err := getPrice()
	if err != nil {
		panic(err)
	}

	if err := saveValue(price); err != nil {
		panic(err)
	}

	fmt.Printf("USD Price: %v\n", price)
}

func getPrice() (float64, error) {
	client := http.Client{
		Timeout: 200 * time.Millisecond,
	}
	res, err := client.Get("http://localhost:8080/cotacao")
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	response := &PriceResponse{}
	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return 0, err
	}

	return response.Price, nil
}

func saveValue(value float64) error {
	if err := os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %v", value)), 0644); err != nil {
		return err
	}
	return nil
}
