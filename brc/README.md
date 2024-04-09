# 1BRC: One Billion Row Challenge in Python

No modelo realizando, foi processado 300.000.000 de linhas em 1 minuto e 25 segundos.

## Bibliotecas utilizadas

panadas
tqdm
time
os
numpy

## Resolução

## Maior Desafio

O maior desafio foi a leitura do arquivo, pois o arquivo é muito grande e não é possível carregar todo o arquivo na memória. Para resolver esse problema, tentei converter o arquivo de `.txt` para `.parquet`. O Resultado foi excelente, o arquivo `.parquet` ficou com 719MB e o `.txt` com 3.9GB. A leitura do arquivo `.parquet` foi muito mais rápida, diminuindo o tempo de processamento e memoria utilizada em mais que a metade.

*NOTA: Também tentei utilizar libs como cudf e processar o arquivo por linha, mas não obtive sucesso.*

## Implementação

A implementação foi realizada em Python, utilizando a biblioteca pandas para realizar o processamento do arquivo. O arquivo foi convertido de `.txt` para `.parquet` para diminuir o tempo de leitura do arquivo. Com o tamanho reduzido do arquivo, foi possível realizar a leitura completa e direta, dessa forma parseando o arquivo e realizando o processamento das linhas e os calculos necessários de min, avg e max. Como no exemplo abaixo:

```python
    # Lê o arquivo Parquet usando pandas
    df = pd.read_parquet(filename)
    # Renomeia as colunas para 'station' e 'value'
    df.columns = ['station', 'value']

    # Define as operações de agregação a serem aplicadas aos valores: mínimo, máximo, média e contagem
    aggregation = {
        'value': ['min', 'max', 'mean', 'count']
    }
    # Agrupa os dados por estação e aplica as operações de agregação
    grouped = df.groupby('station').agg(aggregation)
    # Renomeia as colunas resultantes para facilitar o acesso
    grouped.columns = ['min', 'max', 'avg', 'count']

    # Dicionário para armazenar os objetos Measurement para cada estação
    measurements = {}
    # Itera sobre cada estação, criando um objeto Measurement e armazenando no dicionário
    for id, row in tqdm(grouped.iterrows(), total=grouped.shape[0], desc="Processing"):
        measurements[id] = Measurement(row['min'], row['avg'], row['max'], row['count'])
```

Para poder acompanhar o progresso do processamento, foi utilizado a biblioteca tqdm, que fornece uma barra de progresso para o loop de iteração.

## Google Colab

O Google Colab foi utilizado para realizar o processamento do arquivo, pois o arquivo é muito grande e não é possível carregar todo o arquivo na memória. O Google Colab disponibiliza 12GB de memória RAM e 50GB de armazenamento, o que foi suficiente para realizar o processamento do arquivo.

[Link do Google Colab](https://colab.research.google.com/drive/17klqlnlCFMlXlecbz5dPw1298vHeVeG8?usp=sharing)

## Demonstração

[Acesse o vídeo de demonstração](https://drive.google.com/file/d/1GQnNyypDDCmOr4_QYfHf8qNKJChi6UTQ/view?usp=drive_link)