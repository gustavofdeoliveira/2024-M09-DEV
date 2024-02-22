# Pondera 1 - Simulador de dispositivos IoT

## Objetivo

Criar um simulador de dispositivos IoT utilizando o protocolo MQTT através do uso da biblioteca Eclipse Paho.

## Sensor

Sensor escolhido: SPS30

### Construtor
```go

type SensorData struct {
	PM1         float64 `json:"pm1"`
	PM25        float64 `json:"pm25"`
	PM10        float64 `json:"pm10"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
	Time        string  `json:"time"`
}

```

## Como Rodar

1. Após rodar o mosquito em sua maquina local, com o comando:

```bash

    mosquitto -c mosquitto.conf

``` 
2. Abra um novo terminal e acesse o diretório do projeto, com o seguinte comando:

```bash

    cd P1/src

```
3. Agora, rode o seguinte comando:
    
```bash
    
    chmod +x start.sh
    
```
4. Rode o seguinte comando para instalar as dependências do projeto e executar o simulador:

```bash

    ./start.sh

```

5. O simulador irá rodar e enviará mensagens para o tópico `sensor` a cada 1 segundos.
   


## Demonstração

[Vídeo de monstração](https://drive.google.com/file/d/1-w5XSGHLmXgU7P9mhJWMuavJ7kfO_GqK/view?usp=sharing)