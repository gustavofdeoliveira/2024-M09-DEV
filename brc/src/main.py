import pandas as pd
import time
from tqdm import tqdm


"""# Convert to Parquet"""

def txt_to_parquet(input_file, output_file, separator=','):
    """
    Converte um arquivo TXT para Parquet.

    Args:
    input_file (str): Caminho do arquivo de entrada TXT.
    output_file (str): Caminho do arquivo de saída Parquet.
    separator (str, opcional): Separador utilizado no arquivo TXT. Padrão é ','.
    """
    print("Lendo arquivo .TXT: ")
    df = pd.read_csv(input_file, sep=separator)
    print("Convertendo para .PARQUET:")
    df.to_parquet(output_file, engine='pyarrow')

input_file = '/content/measurements.txt'
output_file = '/content/measurements.parquet'
separator = ';'

txt_to_parquet(input_file, output_file, separator)

"""# Resolve The Challenge"""

import pandas as pd
import numpy as np
import time
from tqdm import tqdm

class Measurement:
    # Construtor da classe Measurement para armazenar as estatísticas de cada estação
    def __init__(self, min_val, avg_val, max_val, count):
        self.min = min_val  # Valor mínimo
        self.avg = avg_val  # Valor médio
        self.max = max_val  # Valor máximo
        self.count = count  # Contagem de medições

def main():
    start_time = time.time()  # Marca o início da execução

    path_file = '/content/measurements.parquet'  # Caminho do arquivo Parquet
    measurements = process_file_with_pandas(path_file)  # Processa o arquivo e calcula estatísticas

    # Exibe os resultados formatados
    print("Resultados:")
    for id, measurement in sorted(measurements.items()):
        print_result(id, measurement)

    # Exibe o tempo total de execução
    print(f"\nTempo de execução: {time.time() - start_time:.2f} segundos.")

def print_result(id, measurement):
    # Formata e exibe os resultados para cada estação
    print(f"{id}={round(measurement.min, 1)}/{round(measurement.avg, 1)}/{round(measurement.max, 1)}")

def process_file_with_pandas(filename):
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

    return measurements

if __name__ == "__main__":
    main()