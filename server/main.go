package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ExternalUSDPriceResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type GetUSDPriceResponse struct {
	Price float64 `json:"price"`
}

func main() {
    http.HandleFunc("GET /cotacao", handleGetPrice)
	http.ListenAndServe(":8080", nil)
}


func handleGetPrice(w http.ResponseWriter, r *http.Request) {
	price, err := getUSDPrice()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendError(w, http.StatusInternalServerError, `{"message":"error getting USD price"}`)
		return
	}

	bid, err := strconv.ParseFloat(price.USDBRL.Bid, 64)
	if err != nil {
		log.Default().Println(err)
		sendError(w, http.StatusInternalServerError, `{"message":"error parsing USD price"}`)
		return
	}

	if err := savePrice(bid); err != nil {
		log.Default().Println(err)
		sendError(w, http.StatusInternalServerError, `{"message":"error saving USD price"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetUSDPriceResponse{Price: bid})
}

func getUSDPrice() (ExternalUSDPriceResponse, error) {
	client := &http.Client{
		Timeout: 200 * time.Millisecond,
	}
	res, err := client.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return ExternalUSDPriceResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ExternalUSDPriceResponse{}, err
	}

	priceResponse := ExternalUSDPriceResponse{}
	err = json.Unmarshal(body, &priceResponse)
	if err != nil {
		return ExternalUSDPriceResponse{}, err
	}

	return priceResponse, nil
}

func savePrice(price float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS prices (id INTEGER PRIMARY KEY, price REAL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO prices(price) VALUES(?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(price)
	if err != nil {
		return err
	}
	return nil
}

func sendError(w http.ResponseWriter, status int, errorBody string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(errorBody))
}

