package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Cotacao struct {
	Bid string `json:bid`
}

func main() {
	http.HandleFunc("/cotacao", func(writer http.ResponseWriter, request *http.Request) {
		cotacao := Cotacao{Bid: "2.81"}

		writer.Header().Set("Content-type", "application/json")
		json.NewEncoder(writer).Encode(cotacao)
	})
	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}
