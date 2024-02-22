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

func publishObject(newObject []map[string]interface{}, singletonClient *MQTTSubscriber) {
	jsonData, err := json.Marshal(newObject)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	token := singletonClient.client.Publish("topic/publisher", 0, false, jsonData)
	token.Wait()
	fmt.Println("Publicado:", string(jsonData))
}

func public(opts *MQTTSubscriber) {
	
	var file = openFile("data.json")
	var bytes = readFile(file)

	var result []map[string]interface{}
	var err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Fatalf("Erro ao decodificar o JSON: %s", err)
	}

	newObject := createObject(result)
	singletonClient := opts

	publishObject(newObject, singletonClient)

	
}
