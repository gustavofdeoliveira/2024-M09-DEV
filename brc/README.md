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

## Implementação - 300.000.000 de linhas

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

## Implementação - 1.000.000.000 de linhas

A implementação foi realizada em Python, utilizando a biblioteca pyarrow para realizar o processamento do arquivo. O arquivo foi convertido de `.txt` para `.parquet` alterando o exemplo fornecido no challenge, para diminuir o tempo de leitura do arquivo e custo de processamento. Com o tamanho reduzido do arquivo, foi possível realizar a leitura linha por linha, dessa forma parseando o arquivo e realizando o processamento das linhas e os calculos necessários de min, avg e max. Como no exemplo abaixo:


```python
def process_file_with_pyarrow(filename):
    parquet_file = pq.ParquetFile(filename)

    # Dicionário para armazenar os objetos Measurement para cada estação
    measurements = {}

    # Iterar através de cada RowGroup no arquivo Parquet
    for i in tqdm(range(parquet_file.num_row_groups), desc="Processing Row Groups"):
        # Lê o RowGroup atual
        table = parquet_file.read_row_group(i)
        df = table.to_pandas()

        df.columns = ['station', 'value']
        aggregation = {
            'value': ['min', 'max', 'mean', 'count']
        }
        grouped = df.groupby('station').agg(aggregation)
        grouped.columns = ['min', 'max', 'avg', 'count']

        # Atualiza o dicionário de medições
        for id, row in grouped.iterrows():
            if id in measurements:
                meas = measurements[id]
                meas.min = min(meas.min, row['min'])
                meas.max = max(meas.max, row['max'])
                meas.avg = (meas.avg * meas.count + row['avg'] * row['count']) / (meas.count + row['count'])
                meas.count += row['count']
            else:
                measurements[id] = Measurement(row['min'], row['avg'], row['max'], row['count'])

    return measurements
```

Para poder acompanhar o progresso do processamento, foi utilizado a biblioteca tqdm, que fornece uma barra de progresso para o loop de iteração.

## Relação ao Projeto

O modúlo atual serviria para processar o arquivo de dados do projeto, realizando o cálculo de min, avg e max para cada estação e armazenando em um dicionário. O arquivo ponderia ser lido e processado por linhas, enviado pro kafka e salvo no mongodb. Posteriomente, selecionar na base conforme as estações e realizar o calculo de min, avg e max.

## Google Colab

O Google Colab foi utilizado para realizar o processamento do arquivo, pois o arquivo é muito grande e não é possível carregar todo o arquivo na memória. O Google Colab disponibiliza 12GB de memória RAM e 50GB de armazenamento, o que foi suficiente para realizar o processamento do arquivo.

### Versão do Google Colab - 1.000.000.000

[Link do Google Colab](https://colab.research.google.com/drive/17klqlnlCFMlXlecbz5dPw1298vHeVeG8?usp=sharing)

## Demonstração

### 300.000.000 de linhas em 1 minuto e 25 segundos

[Acesse o vídeo de demonstração](https://drive.google.com/file/d/1GQnNyypDDCmOr4_QYfHf8qNKJChi6UTQ/view?usp=drive_link)

### 1.000.000.000 5 minutos e 36 segundos

#### Gerando o arquivo

![1712673043607](image/README/create-file.png)

#### Processando o arquivo

![1712673097070](image/README/processed.png)