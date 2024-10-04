package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Data struct {
	Longitude string
	Latitude  string
	Level     string
}

func main() {

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
			log.Println("Método não suportado")
			return
		}

		// Lê o corpo da requisição
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Erro ao ler a requisição", http.StatusInternalServerError)
			log.Println("Erro ao ler a requisição")
			return
		}
		defer req.Body.Close()

		// Decodifica os dados do JSON
		var data Data
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
			log.Println("Erro ao decodificar JSON")
			return
		}
		log.Println("receiving data:", data.Latitude, data.Longitude, data.Level)
	}

	http.Handle("/receive", http.HandlerFunc(helloHandler))

	log.Fatalln(http.ListenAndServe(":8080", nil))

}
