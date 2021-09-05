package main

import (
	"fmt"
	"log"
	"order/queue"
	"os"
	"time"

	"github.com/joho/godotenv"
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


var baseUrlProducts string


func init() {
	// Carregando arquivo .env da pasta product
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Falha ao carregar arquivo .env")
	}

	// Recebe o valor atríbuido a variável PRODUCT_URL=VALOR no arquivo .env
	baseUrlProducts = os.Getenv("PRODUCT_URL")
}


func main() {
	in := make(chan []byte)

	connection := queue.Connect()
	queue.Consumer(connection, in)

	for payload := range in {
		fmt.Println("Mensagem consumida: " + string(payload))
	}
}
