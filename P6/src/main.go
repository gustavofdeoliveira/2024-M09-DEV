package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	godotenv "github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// MQTTSubscriber é uma estrutura que representa um assinante MQTT.
type MQTTSubscriber struct {
	client MQTT.Client
}

// MessageReceiver é uma interface que define um método para receber mensagens MQTT.
type MessageReceiver interface {
	ReceiveMessage(client MQTT.Client, msg MQTT.Message)
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

func createObject(result []map[string]interface{}) []map[string]interface{} {
	var newObject []map[string]interface{}
	for _, item := range result {
		newItem := make(map[string]interface{})
		for key, value := range item {
			// Checa se a chave é "Datetime". Se sim, atualiza para o datetime atual.
			if key == "Datetime" {
				newItem[key] = time.Now().Format(time.RFC3339)
			} else {
				// Para outras chaves, tenta realizar uma operação específica baseada no tipo do valor
				switch v := value.(type) {
				case float64:
					// Se o valor for float64, multiplica por um número aleatório
					newItem[key] = v * rand.Float64()
				default:
					// Para todos os outros tipos, apenas copia o valor
					newItem[key] = value
				}
			}
		}
		newObject = append(newObject, newItem)
	}

	return newObject
}

func publishObject(newObject []map[string]interface{}, singletonClient *MQTTSubscriber) string {
	jsonData, err := json.Marshal(newObject)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ""
	}
	token := singletonClient.client.Publish("topic/publisher", 0, false, jsonData)
	token.Wait()
	fmt.Println("Publicado:", string(jsonData))
	return string(jsonData)
}

var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	fmt.Println("Connected")
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

// NewMQTTSubscriber cria e retorna um novo assinante MQTT.
func NewMQTTSubscriber() *MQTTSubscriber {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
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
		log.Fatalf("Error connecting to MQTT broker: %s", token.Error())
	}

	return &MQTTSubscriber{client: client}
}

// ReceiveMessage implementa o método da interface MessageReceiver para receber mensagens MQTT.
func (s *MQTTSubscriber) ReceiveMessage(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Recebido: %s do tópico: %s\n", msg.Payload(), msg.Topic())
	kafka_producer(msg)
}

func kafka_producer(msg MQTT.Message) {

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS"),
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     os.Getenv("SASL_USERNAME"),
		"sasl.password":     os.Getenv("SASL_PASSWORD"),
	}

	producer, err := kafka.NewProducer(kafkaConf)
	if err != nil {
		log.Fatalf("Falha ao criar produtor: %v", err)
	}
	defer producer.Close()

	topic := os.Getenv("KAFKA_TOPIC")
	fmt.Printf("Conectado ao tópico %s...\n", topic)

	message, err := json.Marshal(msg.Payload())
	if err != nil {
		log.Printf("Erro ao parsear a mensagem: %v", err)
	}

	fmt.Printf("Parsed message: %s\n", message)

	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)

}
func main() {
	subscriber := NewMQTTSubscriber()

	// Assumindo que você corrigiu o método `ReceiveMessage` para ser utilizado corretamente aqui.
	subscriber.client.Subscribe("topic/publisher", 1, func(client MQTT.Client, msg MQTT.Message) {
		// Aqui você chama o método diretamente do seu subscriber.
		subscriber.ReceiveMessage(client, msg)
	})
	var file = readFile(openFile("./data/data.json"))
	result := []map[string]interface{}{}
	var err = json.Unmarshal(file, &result)
	if err != nil {
		log.Fatalf("Erro ao decodificar o JSON: %s", err)
	}
	publishObject(createObject(result), subscriber)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	fmt.Println("Encerrando o programa.")
	subscriber.client.Disconnect(250)
}
