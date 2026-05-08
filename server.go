package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Cotacao struct {
	Bid string `json:bid`
}

type RetornoApiExterna struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", handleCotacao)
	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}

func handleCotacao(writer http.ResponseWriter, reader *http.Request) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Fprint(os.Stderr, "Erro ao fazer requisição: %v\n", err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprint(os.Stderr, "Erro ao ler resposta: %v\n", err)
	}
	var data RetornoApiExterna
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprint(os.Stderr, "Erro ao fazer o parse resposta: %v\n", err)
	}

	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(CotacaoResponse{Bid: data.USDBRL.Bid})
}
