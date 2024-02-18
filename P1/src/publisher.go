package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// SensorData represents the data from the SPS30 sensor
type SensorData struct {
	PM1         float64 `json:"pm1"`
	PM25        float64 `json:"pm25"`
	PM10        float64 `json:"pm10"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
	Time        string  `json:"time"`
}

func public(opts *MQTTSubscriber) {
	singletonClient := opts

	rand.Seed(time.Now().UnixNano())
	sensorData := SensorData{
		PM1:         rand.Float64() * 100,
		PM25:        rand.Float64() * 100,
		PM10:        rand.Float64() * 100,
		Temperature: rand.Float64() * 50,  // Temperature in Celsius
		Humidity:    rand.Float64() * 100, // Relative humidity in percentage
		Pressure:    rand.Float64() * 2000,
		Time:        string(time.Now().Format(time.RFC3339)),
	}
	jsonData, err := json.Marshal(sensorData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Publicar usando a instância única do cliente MQTT.
	token := singletonClient.client.Publish("topic/publisher", 0, false, jsonData)
	token.Wait()
	fmt.Println("Publicado:", string(jsonData))
}
