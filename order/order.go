package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"order/db"
	"order/queue"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nu7hatch/gouuid"
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


func getProductById(id string) Product {

	// Enviando requisição para microsserviço de produto
	response, err := http.Get(baseUrlProducts + "/products/" + id)
	if err != nil {
		fmt.Printf("Falha ao carregar requisição HTTP %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)

	var product Product
	json.Unmarshal(data, &product)

	return product
}


func CreateOrder(payload []byte) {
	var order Order
	json.Unmarshal(payload, &order)

	uuid, _ := uuid.NewV4()

	order.Uuid = uuid.String()
	order.Status = "Pendente"
	order.CreatedAt = time.Now()

	saveOrder(order)
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


func main() {
	in := make(chan []byte)

	connection := queue.Connect()
	queue.Consumer(connection, in)

	for payload := range in {
		fmt.Println("Mensagem consumida: " + string(payload))

		// Criar ordem de serviço
		CreateOrder(payload)
	}
}
