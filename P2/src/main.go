package main

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// MessageReceiver é uma interface que define um método para receber mensagens MQTT.
type MessageReceiver interface {
	ReceiveMessage(client MQTT.Client, msg MQTT.Message)
}

// MQTTSubscriber é uma estrutura que representa um assinante MQTT.
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

func public(opts *MQTTSubscriber) string{

	var file = openFile("data.json")
	var bytes = readFile(file)

	var result []map[string]interface{}
	var err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Fatalf("Erro ao decodificar o JSON: %s", err)
	}

	newObject := createObject(result)
	singletonClient := opts

	var objectPublicated = publishObject(newObject, singletonClient)
	return objectPublicated

}

// NewMQTTSubscriber cria e retorna um novo assinante MQTT.
func NewMQTTSubscriber(messageReceiver MessageReceiver) *MQTTSubscriber {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetClientID("go_subscriber")
	
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	client.Subscribe("topic/publisher", 1, func(client MQTT.Client, msg MQTT.Message) {
		messageReceiver.ReceiveMessage(client, msg)
	})

	return &MQTTSubscriber{client: client}
}

// ReceiveMessage implementa o método da interface MessageReceiver para receber mensagens MQTT.
func (s *MQTTSubscriber) ReceiveMessage(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Recebido: %s do tópico: %s\n", msg.Payload(), msg.Topic())
}

func sub() {
	subscriber := NewMQTTSubscriber(&MQTTSubscriber{})

	// Cria um canal para capturar sinais do sistema operacional (Ctrl+C)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Aguarda sinal de término
	<-sigCh
	fmt.Println("Encerrando o programa.")
	subscriber.client.Disconnect(250)
}

func main() {
	public(NewMQTTSubscriber(&MQTTSubscriber{}))
	sub()
}