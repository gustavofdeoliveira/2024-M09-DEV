# Integração com o Postgress

## Requisitos

1. Ter o Docker instalado

## Como rodar

1. Acesse a pasta `P5/src/docker` e execute o comando `docker-compose up -d` ou `docker compose up -d`
2. Espere alguns segundos para que o postgres, metabase e mongo sejam inicializados
3. Acesse o metabase em `http://localhost:3000` e configure a conexão com o banco de dados postgres.

## Configuração do Metabase e Postgres

1. pós a etapa de configuração do metabase, você pode configurar a integração com o banco de dados postgres acessando a aba `Databases` e clicando em `Add Database` e selecionando `Postgres` e preenchendo os campos com as informações do banco de dados.

### Informações do banco de dados

```env
POSTGRES_PASSWORD=minhaSenhaSegura
POSTGRES_USER=meuUsuario
POSTGRES_DB=meuBancoDeDados
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```

NOTA: Preencha os campos com as informações do acima.

## Verificando Integração

1. Após a configuração do banco de dados, você pode parar o container do metabase e do postgres e iniciar novamente com o comando `docker-compose up -d` ou `docker compose up -d` para verificar se a integração foi realizada com sucesso. Se as informações de configuração estiverem corretas, o metabase irá se conectar ao banco de dados postgres e suas configurações estarão salvas.
2. Acesse a pasta `P5/src/docker` e verique se as pastas `db-data`, `metabase-data` e `mongo-data` foram criadas.

## Demonstração

[Acesse o vídeo de demonstração](https://drive.google.com/file/d/1W-fTfVh-jr9AndRh8cs00RPn1d367xOd/view?usp=drive_link)