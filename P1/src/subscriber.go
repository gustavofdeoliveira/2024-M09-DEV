package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MessageReceiver é uma interface que define um método para receber mensagens MQTT.
type MessageReceiver interface {
	ReceiveMessage(client MQTT.Client, msg MQTT.Message)
}

// MQTTSubscriber é uma estrutura que representa um assinante MQTT.
type MQTTSubscriber struct {
	client MQTT.Client
}

// NewMQTTSubscriber cria e retorna um novo assinante MQTT.
func NewMQTTSubscriber(opts *MQTT.ClientOptions, messageReceiver MessageReceiver) *MQTTSubscriber {
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

func main() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetClientID("go_subscriber")

	subscriber := NewMQTTSubscriber(opts, &MQTTSubscriber{})

	// Cria um canal para capturar sinais do sistema operacional (Ctrl+C)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Loop principal
	go func() {
		for {
			if subscriber.client.IsConnected() {
				public(subscriber) // Assumindo que `public` agora aceita um cliente como argumento
			} else {
				fmt.Println("Cliente MQTT não conectado. Tentando reconectar...")
				subscriber.client.Connect()
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Aguarda sinal de término
	<-sigCh
	fmt.Println("Encerrando o programa.")
	subscriber.client.Disconnect(250)
}
