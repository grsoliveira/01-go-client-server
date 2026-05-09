package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

	content := fmt.Sprintf("Dólar: %s", cotacao.Bid)

	err = os.WriteFile("cotacao.txt", []byte(content), 0644)
	if err != nil {
		log.Println("Erro na escrita do arquivo", err)
		return
	}

	fmt.Println("Cotação salva com sucesso!")
}
