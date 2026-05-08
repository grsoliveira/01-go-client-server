package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/cotacao", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Servidor rodando!")
	})
	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}
