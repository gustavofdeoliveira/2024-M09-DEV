package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	time "time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	godotenv "github.com/joho/godotenv"
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

func TestConnection(t *testing.T) {
	fmt.Println("TestConnection")
	subscriber := NewMQTTSubscriber()

	if subscriber.client.IsConnected() {
		fmt.Println("Conectado")
	} else {
		t.Errorf("Erro de conexão")
	}
	subscriber.client.Disconnect(250)
}

func TestBrokerPublicAndRecevedMessage(t *testing.T) {
	fmt.Println("TestBrokerPublicAndRecevedMessage")
	var file = openFile("./data/data.json")
	var bytes = readFile(file)

	var result []map[string]interface{}
	var err = json.Unmarshal(bytes, &result)
	if err != nil {
		t.Fatalf("Erro ao decodificar o JSON: %s", err)
	}
	var subscriber = NewMQTTSubscriber()
	for _, item := range result {
		newObject := createObject(item)
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file: %s", err)
		}
		time.Sleep(3 * time.Second)
		subscriber.client.Subscribe("topic/publisher", 1, func(client MQTT.Client, msg MQTT.Message) {
			subscriber.ReceiveMessage(client, msg)
			if msg.Payload() == nil {
				t.Errorf("Mensagem vazia")
			}
		})
		publishObject(newObject, subscriber)
	}
	subscriber.client.Disconnect(250)
}

func TestKafkaPublicAndRecevedMessage(t *testing.T) {
	fmt.Println("TestKafkaPublicAndRecevedMessage")
	var file = openFile("./data/data.json")
	var bytes = readFile(file)

	var result []map[string]interface{}
	var err = json.Unmarshal(bytes, &result)
	if err != nil {
		t.Fatalf("Erro ao decodificar o JSON: %s", err)
	}
	var subscriber = NewMQTTSubscriber()
	for _, item := range result {
		newObject := createObject(item)
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file: %s", err)
		}
		time.Sleep(3 * time.Second)

		var publishMessage string 

		subscriber.client.Subscribe("topic/publisher", 1, func(client MQTT.Client, msg MQTT.Message) {
			subscriber.ReceiveMessage(client, msg)
			fmt.Printf("[BROKER] Recebido: %s do tópico: %s\n", msg.Payload(), msg.Topic())
			consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
				"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS"),
				"group.id":          "go-kafka-consumer",
				"security.protocol": "SASL_SSL",
				"sasl.mechanisms":   "PLAIN",
				"sasl.username":     os.Getenv("SASL_USERNAME"),
				"sasl.password":     os.Getenv("SASL_PASSWORD"),
			})
			if err != nil {
				log.Fatalf("[CONSUMER] Falha ao criar produtor: %v", err)
			}
			defer consumer.Close()

			topic := os.Getenv("KAFKA_TOPIC")
			fmt.Printf("[CONSUMER] Conectado ao tópico %s...\n", topic)

			consumer.SubscribeTopics([]string{topic}, nil)
			for {
				message, err := consumer.ReadMessage(-1)
				if err == nil {
					if string(message.Value) != publishMessage {
						log.Fatalf("[CONSUMER] Mensagem não recebida: %s", message.Value)
					}
				} else {
					fmt.Printf("[CONSUMER] error: %v (%v)\n", err, message)
					break
				}
			}
			consumer.Close()
		})
		publishMessage = publishObject(newObject, subscriber)
	}
	subscriber.client.Disconnect(250)
}

func TestCreateAndPublisObject(t *testing.T) {
	fmt.Println("TestCreateAndPublisObject")
	var result []map[string]interface{}
	bytes := []byte(`[{"Datetime":"2021-09-01T12:00:00Z","Value":10.0}]`)
	json.Unmarshal(bytes, &result)
	subscriber := NewMQTTSubscriber()
	for _, item := range result {
		newObject := createObject(item)

		if newObject == nil {
			t.Errorf("Erro ao criar objeto")
		}
		publishObject(newObject, subscriber)
	}
	subscriber.client.Disconnect(250)
}
