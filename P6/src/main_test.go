package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	godotenv "github.com/joho/godotenv"
	"log"
	"os"
	"testing"
	time "time"
)

func TestOpenFileSuccess(t *testing.T) {
	fmt.Println("TestOpenFileSuccess")

	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Erro ao criar arquivo temporário: %s", err)
	}
	tmpfilePath := tmpfile.Name()

	defer os.Remove(tmpfilePath)
	tmpfile.Close()

	file := openFile(tmpfilePath)
	if file == nil {
		t.Errorf("openFile retornou nil para um arquivo existente")
	}
	file.Close()
}

func TestReadFileSuccess(t *testing.T) {
	fmt.Println("TestReadFileSuccess")
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Erro ao criar arquivo temporário: %s", err)
	}
	tmpfilePath := tmpfile.Name()

	// Cleanup: Garante que o arquivo temporário seja removido após o teste.
	defer os.Remove(tmpfilePath)
	tmpfile.Close()

	// Teste: Tenta abrir o arquivo temporário.
	file := openFile(tmpfilePath)
	if file == nil {
		t.Errorf("openFile retornou nil para um arquivo existente")
	}
	bytes := readFile(file)

	if bytes == nil {
		t.Errorf("readFile retornou nil para um arquivo existente")
	}
	file.Close()

}

func TestCreateAndPublisObject(t *testing.T) {
	fmt.Println("TestCreateAndPublisObject")
	var result []map[string]interface{}
	bytes := []byte(`[{"Datetime":"2021-09-01T12:00:00Z","Value":10.0}]`)
	json.Unmarshal(bytes, &result)
	for _, item := range result {
		newObject := createObject(item)

		if newObject == nil {
			t.Errorf("Erro ao criar objeto")
		}

		subscriber := NewMQTTSubscriber()

		publishObject(newObject, subscriber)
	}

}

func TestPublicAndRecevedMessage(t *testing.T) {
	fmt.Println("TestPublicAndRecevedMessage")
	var file = openFile("./data/data.json")
	var bytes = readFile(file)

	var result []map[string]interface{}
	var err = json.Unmarshal(bytes, &result)
	if err != nil {
		t.Fatalf("Erro ao decodificar o JSON: %s", err)
	}
	for _, item := range result {
		newObject := createObject(item)

		var subscriber = NewMQTTSubscriber()
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file: %s", err)
		}
		time.Sleep(3 * time.Second)
		subscriber.client.Subscribe("topic/publisher", 1, func(client MQTT.Client, msg MQTT.Message) {
			subscriber.ReceiveMessage(client, msg)
			producer, err := kafka.NewProducer(&kafka.ConfigMap{
				"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS"),
				"client.id":         "go-kafka-producer",
				"security.protocol": "SASL_SSL",
				"sasl.mechanisms":   "PLAIN",
				"sasl.username":     os.Getenv("SASL_USERNAME"),
				"sasl.password":     os.Getenv("SASL_PASSWORD"),
			})
			if err != nil {
				log.Fatalf("Falha ao criar produtor: %v", err)
			}
			producer.Flush(15 * 1000)
			consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
				"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS"),
				"sasl.username":     os.Getenv("SASL_USERNAME"),
				"sasl.password":     os.Getenv("SASL_PASSWORD"),
				"security.protocol": "SASL_SSL",
				"sasl.mechanisms":   "PLAIN",
				"group.id":          "go-consumer-group",
				"auto.offset.reset": "earliest",
			})

			if err != nil {
				panic(err)
			}
			defer consumer.Close()

			topic := os.Getenv("KAFKA_TOPIC")
			fmt.Printf("Conectado ao tópico %s...\n", topic)

			consumer.SubscribeTopics([]string{topic}, nil)

			for {
				msg, err := consumer.ReadMessage(-1)
				if err == nil {
					fmt.Printf("Consumer received message: %s\n", string(msg.Value))
				} else {
					fmt.Printf("Consumer error: %v (%v)\n", err, msg)
					break
				}
			}
		})

		publishObject(newObject, subscriber)

	}
}

// func TestConnection(t *testing.T) {

// 	subscriber := NewMQTTSubscriber()

// 	if subscriber.client.IsConnected() {
// 		fmt.Println("Conectado")
// 	} else {
// 		t.Errorf("Erro de conexão")
// 	}
// 	subscriber.client.Disconnect(250)
// }
