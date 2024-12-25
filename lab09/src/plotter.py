import os
import pandas as pd
import matplotlib.pyplot as plt


def parse_duration(duration_str):
    """Парсит строку времени в миллисекунды"""
    if 'ms' in duration_str:
        return float(duration_str.replace('ms', ''))
    elif 'µs' in duration_str:
        return float(duration_str.replace('µs', '')) / 1000
    elif 's' in duration_str:
        return float(duration_str.replace('s', '')) * 1000
    return 0


def plot_scenario_data(csv_file, output_dir):
    """Строит график для одного сценария на основе .csv файла"""
    data = pd.read_csv(csv_file)

    data['DBQueryTime'] = data['DBQueryTime'].apply(parse_duration)
    data['CacheQueryTime'] = data['CacheQueryTime'].apply(parse_duration)

    data.sort_values(by='Second', inplace=True)

    plt.figure(figsize=(10, 6))
    plt.plot(data['Second'], data['DBQueryTime'], label='DBQueryTime (ms)', marker='o')
    plt.plot(data['Second'], data['CacheQueryTime'], label='CacheQueryTime (ms)', marker='x')
    plt.xlabel('Time of Scenario (Second)')
    plt.ylabel('Query Time (ms)')
    plt.title(f'Query Time Comparison - {os.path.basename(csv_file)}')
    plt.legend()
    plt.grid(True)

    output_file = os.path.join(output_dir, os.path.basename(csv_file).replace('.csv', '_plot.png'))
    plt.savefig(output_file)
    plt.close()
    print(f"График сохранен: {output_file}")


def main():
    input_dir = os.path.join('.', 'data')
    output_dir = os.path.join('.', 'data', 'plots')

    os.makedirs(output_dir, exist_ok=True)

    csv_files = [os.path.join(input_dir, f) for f in os.listdir(input_dir) if f.endswith('.csv')]

    if not csv_files:
        print(f"В папке {input_dir} нет .csv файлов")
        return

    for csv_file in csv_files:
        print(f"Обрабатывается файл: {csv_file}")
        plot_scenario_data(csv_file, output_dir)


if __name__ == "__main__":
    main()
