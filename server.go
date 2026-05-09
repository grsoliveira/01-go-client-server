package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type RetornoApiExterna struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	createTable(db)

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		handleCotacao(w, db)
	})
	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}

func createTable(db *sql.DB) {
	query := "CREATE TABLE IF NOT EXISTS cotacoes ( " +
		"id INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"bid TEXT, " +
		"created_at DATETIME DEFAULT CURRENT_TIMESTAMP " +
		");"

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func handleCotacao(w http.ResponseWriter, db *sql.DB) {
	ctxAPI, cancelAPI := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancelAPI()

	req, err := http.NewRequestWithContext(
		ctxAPI,
		"GET",
		"https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar requisição: %v\n", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(os.Stderr, "Erro ao ler resposta: %v\n", err)
	}

	var data RetornoApiExterna
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprint(os.Stderr, "Erro ao fazer o parse resposta: %v\n", err)
	}

	cotacao := Cotacao{
		Bid: data.USDBRL.Bid,
	}

	salvarCotacao(db, cotacao.Bid)

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(cotacao)
}

func salvarCotacao(db *sql.DB, Bid string) {
	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDB()

	query := "INSERT INTO cotacoes(bid) VALUES (?)"

	_, err := db.ExecContext(ctxDB, query, Bid)
	if err != nil {
		log.Println("Erro ao salvar cotação no banco de dados", err)
	}
}
