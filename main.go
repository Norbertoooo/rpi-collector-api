package main

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"io"
	"log"
	"net/http"
	"time"
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

		connection, _ := amqp.Dial("amqps://gktopult:0SipH52ziVwwBCICELz96krbDbgvAYKs@jackal.rmq.cloudamqp.com/gktopult")
		defer connection.Close()
		channel, _ := connection.Channel()
		defer channel.Close()

		msg := amqp.Publishing{
			DeliveryMode: 1,
			Timestamp:    time.Now(),
			ContentType:  "application/json",
			Body:         body,
		}

		err = channel.Publish("receive.data.exchange", "ping", false, false, msg)
		if err != nil {
			log.Fatalln("erro ao publicar mensagem na fila")
			return
		}

	}

	connection, err := amqp.Dial("amqps://gktopult:0SipH52ziVwwBCICELz96krbDbgvAYKs@jackal.rmq.cloudamqp.com/gktopult")

	if err != nil {
		log.Fatalln("erro ao conectar com rabbit")
	}
	defer connection.Close()

	channel, _ := connection.Channel()
	defer channel.Close()

	q, _ := channel.QueueDeclare("receive.data.queue", true, false, false, true, nil)
	err = channel.QueueBind(q.Name, "#", "receive.data.exchange", false, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	http.Handle("/receive", http.HandlerFunc(helloHandler))

	log.Fatalln(http.ListenAndServe(":8080", nil))

}
