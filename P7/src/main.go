package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	godotenv "github.com/joho/godotenv"
)

type MQTTSubscriber struct {
	client MQTT.Client
}

func openFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo: %s", err)
	}
	return file
}

func readFile(file *os.File) []byte {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo: %s", err)
	}
	return bytes
}

func createObject(result map[string]interface{}) map[string]interface{} {
	newItem := make(map[string]interface{})
	for key, value := range result {
		if key == "Datetime" {
			newItem[key] = time.Now().Format(time.RFC3339)
		} else {
			switch v := value.(type) {
			case float64:
				newItem[key] = v * rand.Float64()
			default:
				newItem[key] = value
			}
		}
	}
	return newItem
}

func publishObject(newObject map[string]interface{}, singletonClient *MQTTSubscriber) string {
	jsonData, err := json.Marshal(newObject)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ""
	}
	token := singletonClient.client.Publish("ponderada", 0, false, jsonData)
	token.Wait()
	fmt.Println("[BROKER] Publicado:", string(jsonData))
	return string(jsonData)
}

var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	fmt.Println("[BROKER] Connected")
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	fmt.Printf("[BROKER] Connection lost: %v", err)
}

func NewMQTTSubscriber() *MQTTSubscriber {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("[ENV] Error loading: %s", err)
	}

	var broker = os.Getenv("BROKER_ADDR")
	var port = 8883
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%d/mqtt", broker, port))
	opts.SetUsername(os.Getenv("HIVE_USER"))
	opts.SetPassword(os.Getenv("HIVE_PSWD"))
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	opts.SetClientID("go_subscriber")

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("[BROKER] Error connecting to MQTT: %s", token.Error())
	}

	return &MQTTSubscriber{client: client}
}

func (s *MQTTSubscriber) ReceiveMessage(client MQTT.Client, msg MQTT.Message) {
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
			fmt.Printf("[CONSUMER] Received message: %s\n", string(message.Value))
			
		} else {
			fmt.Printf("[CONSUMER] error: %v (%v)\n", err, message)
			break
		}
	}
	consumer.Close()

}

func main() {
	subscriber := NewMQTTSubscriber()

	subscriber.client.Subscribe("ponderada", 1, func(client MQTT.Client, msg MQTT.Message) {
		subscriber.ReceiveMessage(client, msg)
	})
	var file = readFile(openFile("./data/data.json"))
	result := []map[string]interface{}{}
	var err = json.Unmarshal(file, &result)
	if err != nil {
		log.Fatalf("Erro ao decodificar o JSON: %s", err)
	}
	for _, item := range result {
		publishObject(createObject(item), subscriber)
		time.Sleep(5 * time.Second)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	fmt.Println("Encerrando o programa.")
	subscriber.client.Disconnect(250)
}
