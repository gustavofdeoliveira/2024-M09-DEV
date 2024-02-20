package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

// SensorData represents the data from the SPS30 sensor

func public(opts *MQTTSubscriber) {
	singletonClient := opts

	var result []map[string]interface{}

	file, err := os.Open("./data.json")
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo: %s", err)
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo: %s", err)
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Fatalf("Erro ao decodificar o JSON: %s", err)
	}
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

	 jsonData, err := json.Marshal(newObject)
	 if err != nil {
	 	fmt.Println("Error marshalling JSON:", err)
	 	return
	 }

	 // Publicar usando a instância única do cliente MQTT.
	 token := singletonClient.client.Publish("topic/publisher", 0, false, jsonData)
	 token.Wait()
	 fmt.Println("Publicado:", string(jsonData))
}
