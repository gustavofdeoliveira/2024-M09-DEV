package main

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"testing"
)

func TestOpenFileSuccess(t *testing.T) {
	fmt.Println("TestOpenFileSuccess")
	// Setup: Cria um arquivo temporário para teste.
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

	// Não esqueça de fechar o arquivo aberto pela função openFile.
	file.Close()
}

func TestReadFileSuccess(t *testing.T) {
	fmt.Println("TestReadFileSuccess")
	// Setup: Cria um arquivo temporário para teste.
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

	newObject := createObject(result)
	if newObject == nil {
		t.Errorf("Erro ao criar objeto")
	}

	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetClientID("go_subscriber")

	subscriber := NewMQTTSubscriber(&MQTTSubscriber{})

	publishObject(newObject, subscriber)

}

func TestPublicAndRecevedMessage(t *testing.T) {
	fmt.Println("TestPublicAndRecevedMessage")
	var file = openFile("data.json")
	var bytes = readFile(file)

	var result []map[string]interface{}
	var err = json.Unmarshal(bytes, &result)
	if err != nil {
		t.Fatalf("Erro ao decodificar o JSON: %s", err)
	}

	newObject := createObject(result)

	var subscriber = NewMQTTSubscriber(&MQTTSubscriber{})

	var jsonObject = publishObject(newObject, subscriber)

	var messageReceiver = &MQTTSubscriber{}

	messageChannel := make(chan string)

	subscriber.client.Subscribe("topic/publisher", 1, func(client MQTT.Client, msg MQTT.Message){
		messageReceiver.ReceiveMessage(client, msg)

		messageChannel <- string(msg.Payload())
	})

	receivedMessage := <-messageChannel

	if receivedMessage != jsonObject{
		t.Errorf("Erro ao receber mensagem")
	
	}

	close(messageChannel)
}

// func TestTimeSend(t *testing.T) {
// 	fmt.Printf("TestTimeSend")
	
//     var result []map[string]interface{}
//     bytes := []byte(`[{"Datetime":"2021-09-01T12:00:00Z","Value":10.0}]`)
//     json.Unmarshal(bytes, &result)

    
   
// 	newObject := createObject(bytes)
// 	if newObject == nil {
// 		t.Errorf("Erro ao criar objeto")
// 	}

// 	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
// 	opts.SetClientID("go_subscriber")

// 	subscriber := NewMQTTSubscriber(opts, &MQTTSubscriber{})

// 	publishObject(newObject, subscriber)

// }

func TestConnection(t *testing.T) {

	subscriber := NewMQTTSubscriber(&MQTTSubscriber{})

	if subscriber.client.IsConnected() {
		fmt.Println("Conectado")
	} else {
		t.Errorf("Erro de conexão")
	}
	subscriber.client.Disconnect(250)
}
