package main

import (
	"encoding/json"
	"fmt"
	"log"
	"order/db"
	"order/queue"
	"time"

	"github.com/joho/godotenv"
	"github.com/nu7hatch/gouuid"
	"github.com/streadway/amqp"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float64 `json:"price,string"`
}

type Order struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	ProductId string    `json:"product_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,string"`
}


func init() {
	// Carregando arquivo .env da pasta product
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Falha ao carregar arquivo .env")
	}
}


func CreateOrder(payload []byte) Order {
	var order Order
	json.Unmarshal(payload, &order)

	uuid, _ := uuid.NewV4()

	order.Uuid = uuid.String()
	order.Status = "Pendente"
	order.CreatedAt = time.Now()

	saveOrder(order)

	return order
}


func saveOrder(order Order) {
	json, _ := json.Marshal(order)

	connection := db.Connect()
	// Salva no Banco de Dados Redis o conjunto: (chave, valor)
	err := connection.Set(order.Uuid, string(json), 0).Err()

	if err != nil {
		panic(err.Error())
	}
}


func notifyOrderCreated(order Order, ch *amqp.Channel) {
	json, _ := json.Marshal(order)
	queue.Publisher(json, "order_ex", "", ch)
}


func main() {
	in := make(chan []byte)

	connection := queue.Connect()
	queue.Consumer("checkout_queue", connection, in)

	for payload := range in {
		fmt.Println("Mensagem consumida: " + string(payload))

		// Criar ordem de serviço
		notifyOrderCreated(CreateOrder(payload), connection)
	}
}
