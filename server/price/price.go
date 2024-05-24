package price

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

type USDPriceResponse struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type PriceUserResponse struct {
	Price float64 `json:"price"`
}

type PriceController struct{}

const USD_PRICE_URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func (p PriceController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	price, err := getUSDPrice()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"error getting USD price"}`))
		return
	}

	log.Default().Println(price)

	bid, err := strconv.ParseFloat(price.USDBRL.Bid, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"error getting USD price"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PriceUserResponse{Price: bid})
}

func getUSDPrice() (USDPriceResponse, error) {
	res, err := http.Get(USD_PRICE_URL)
	if err != nil {
		return USDPriceResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return USDPriceResponse{}, err
	}

	priceResponse := USDPriceResponse{}
	err = json.Unmarshal(body, &priceResponse)
	if err != nil {
		return USDPriceResponse{}, err
	}

	return priceResponse, nil
}
