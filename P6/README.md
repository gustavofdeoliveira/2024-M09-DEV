# Integração com o Postgress

## Requisitos

1. Ter uma conta no HiveMQ Cloud e no Confluent Cloud

## Configuração

Na pasta raiz do projeto, rode o comando `cd src/` para acessar a pasta `src`. E crie um arquivo `.env` com as seguintes informações:

### Informações do banco de dados

```env
HIVE_USER=""
HIVE_PSWD=""
BROKER_ADDR=""
BOOTSTRAP_SERVERS=""
SASL_USERNAME=""
SASL_PASSWORD=""
KAFKA_TOPIC="ponderada"
```

NOTA: Preencha os campos com as informações do acima.

Após a criação do arquivo `.env`, acesse a página do HiveMQ Cloud e crie uma `Integrations` para o confluent, após ter configurado e criado um cluster no Confluent Cloud e tópico denominado `ponderada`. Na integração, copie o `username`, `password`, `bootstrap` e preencha os campos `Source - Topic` e `Destination Topic` respectivamente com o valor `ponderada`.


## Execução
Para executar o projeto, rode os seguintes comandos:

```bash
chmod +x ./start.sh && ./start.sh
```


## Demonstração

[Acesse o vídeo de demonstração](https://drive.google.com/file/d/1PqVF06liLjh5QyC7tkKJdTDLgpcf4cim/view?usp=sharing)