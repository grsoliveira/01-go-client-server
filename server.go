package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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
	http.HandleFunc("/cotacao", handleCotacao)
	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}

func handleCotacao(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(Cotacao{Bid: data.USDBRL.Bid})
}
