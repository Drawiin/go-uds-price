package main

import (
	"database/sql"
	"testing"
)

func TestGetUSDPrice(t *testing.T) {
	price, err := getUSDPrice()
	if err != nil {
		t.Errorf("Erro ao obter o preço do dólar: %v", err)
	}

	if price.USDBRL.Bid == "" {
		t.Errorf("O preço do dólar está vazio")
	}
}

func TestSavePrice(t *testing.T) {
	price := 5.50 // Preço de exemplo
	err := savePrice(price)
	if err != nil {
		t.Errorf("Erro ao salvar o preço: %v", err)
	}

	db, err := sql.Open("sqlite3", "./database.db")

	if err != nil {
		t.Errorf("Erro ao abrir o banco de dados: %v", err)
	}
	defer db.Close()

	var savedPrice float64
	err = db.QueryRow("SELECT price FROM prices").Scan(&savedPrice)
	if err != nil {
		t.Errorf("Erro ao consultar o preço salvo: %v", err)
	}
}