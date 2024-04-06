# Integração do HiveMQ Cloud com Confluent Cloud

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
## Testes

### Testes Unitários
Para testar o projeto, rode o seguinte comando:

```bash
chmod +x ./test.sh && ./test.sh
```
### Testes de Integração

Para testar se a integração está funcionando, acesse o HiveMQ Cloud e conecte-se ao broker, após isso, publique uma mensagem no tópico `ponderada` e verifique se a mensagem foi recebida no Confluent Cloud.

Outro teste que pode ser feito é acessar o Confluent Cloud e verificar se o número de mensagem recebidas no tópico `ponderada` aumentou. Ainda no HiveMQ Cloud, é possível verificar se o número de mensagens enviadas aumentou também.

## Demonstração

[Acesse o vídeo de demonstração](https://drive.google.com/file/d/1PqVF06liLjh5QyC7tkKJdTDLgpcf4cim/view?usp=sharing)