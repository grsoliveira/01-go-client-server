package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CotacaoClient struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"http://localhost:8080/cotacao",
		nil)
	if err != nil {
		log.Println("Erro ao criar requisição de cotações", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro ao realizar a requisição de cotações", err)
		return
	}
	defer resp.Body.Close()

	var cotacao CotacaoClient
	err = json.NewDecoder(resp.Body).Decode(&cotacao)
	if err != nil {
		log.Println("Erro no processamento do retorno da requisição", err)
		return
	}

	fmt.Printf("A cotação retornada foi de %v", cotacao.Bid)
}
