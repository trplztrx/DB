import csv
import random as r
from faker import Faker
from datetime import datetime
import uuid
import os

fake = Faker("ru_RU")

records_num = 1001
user_type = ["физ", "юр"]

def generate_users():
    users = []
    for _ in range(records_num):
        cur_user_type = r.choice(user_type)
        if cur_user_type == "физ":
            name = fake.name()
        else:
            name = fake.company()
        
        user = {
            "user_type": cur_user_type,
            "name": name,
            "phone": fake.unique.phone_number(),
            "email": fake.unique.email()
        }
        users.append(user)
    return users

def generate_file_name(table_name):
    unique_id = uuid.uuid4()
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    return f"{unique_id}_{table_name}_{timestamp}.csv"

def save_to_csv(data, table_name, fieldnames):
    file_name = generate_file_name(table_name)
    directory = "./data"
    os.makedirs(directory, exist_ok=True)
    file_path = os.path.join(directory, file_name)
    with open(file_path, "w", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(data)
    print(f"File saved: {file_path}")

if __name__ == "__main__":
    users = generate_users()
    save_to_csv(users, "users", ["user_type", "name", "phone", "email"])
