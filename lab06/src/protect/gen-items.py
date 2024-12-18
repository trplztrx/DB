import csv
import random
from datetime import datetime

NUM_RECORDS = 800000
CSV_FILE = 'items_data.csv'

def generate_item_name():
    prefixes = ['Super', 'Mega', 'Ultra', 'Pro', 'Eco', 'Smart']
    suffixes = ['Widget', 'Gadget', 'Tool', 'Item', 'Device', 'Object']
    return f"{random.choice(prefixes)} {random.choice(suffixes)}"

def generate_description():
    descriptions = [
        "High-quality product for daily use.",
        "Eco-friendly and sustainable.",
        "Reliable and durable construction.",
        "Lightweight and easy to handle.",
        "Perfect for professionals and beginners.",
        "Comes with a one-year warranty."
    ]
    return random.choice(descriptions)

def generate_price():
    return round(random.uniform(10.00, 1000.00), 2)

def generate_quantity():
    return random.randint(0, 500)

def generate_csv_file():
    with open(CSV_FILE, mode='w', newline='', encoding='utf-8') as file:
        writer = csv.writer(file)
        writer.writerow(['name', 'description', 'quantity', 'price', 'created_at', 'updated_at'])

        for _ in range(NUM_RECORDS):
            name = generate_item_name()
            description = generate_description()
            quantity = generate_quantity()
            price = generate_price()
            created_at = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
            updated_at = created_at

            writer.writerow([name, description, quantity, price, created_at, updated_at])

    print(f"CSV файл с {NUM_RECORDS} записями успешно создан: {CSV_FILE}")

if __name__ == "__main__":
    generate_csv_file()
