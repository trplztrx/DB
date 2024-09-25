import csv
import random as r
from faker import Faker
import re
import itertools
from datetime import timedelta, datetime

fake = Faker("ru_RU")

records_num = 1001
max_item_num = 10
max_delivery_days = 100
max_item_quantity = 100

order_status = ["создан", "отправлен", "доставлен", "отменен"]
user_type = ["физик", "юрик"]
address_type = ["отправитель", "получатель"]
courier_type = ["свободен", "в пути", "занят"]

users = []
for _ in range(records_num):
    cur_user_type = r.choice(user_type)
    if cur_user_type == "физик":
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

couriers = []
for _ in range(records_num): 
    name = fake.name().split()

    courier = {
        "name": name[1],
        "surname": name[0],
        "patronymic": name[2],
        "phone": fake.unique.phone_number(),
        "status": r.choice(courier_type)
    }
    couriers.append(courier)

addresses = []
ids_pickup_address = []
ids_delivery_address = []
for i in range(records_num):
    cur_address_type = r.choice(address_type)
    address = {
        "region": fake.region(),
        "city": fake.city(),
        "address_street": fake.street_address(),
        "address_type": cur_address_type
    }
    if cur_address_type == 'отправитель':
        ids_pickup_address.append(i + 1)
    else:
        ids_delivery_address.append(i + 1)
    
    addresses.append(address)

orders = []
for _ in range (records_num):
    created_at = fake.date_time_between(end_date=f'-{max_delivery_days}d')
    status = r.choice(order_status)
    if status == "создан":
        updated_at = created_at 
    else:
        updated_at = created_at + timedelta(days=fake.random_int(min=1, max=max_delivery_days))
    sender_user_id, receiver_user_id = r.sample(range(1, len(users) + 1), 2)

    order = {
        "created_at": created_at,
        "updated_at": updated_at,
        "status": status,
        "delivery_cost": fake.pydecimal(left_digits=8, right_digits=2, positive=True),
        "sender_user_id": sender_user_id,
        "receiver_user_id": receiver_user_id,
        "pickup_address_id": r.choice(ids_pickup_address),
        "delivery_address_id": r.choice(ids_delivery_address)
    }
    orders.append(order)

order_items = []
for i in range(records_num):
    order_id = i + 1
    for __ in range(r.randint(1, max_item_num)):
        item = {
            "order_id": order_id,
            "item_name": fake.word(),
            "quantity": r.randint(1, max_item_quantity)
        }
        order_items.append(item)




order_statuses = []
for i in range(len(orders)):
    order_created_at = orders[i]["created_at"]
    order_updated_at = orders[i]["updated_at"]
    status = orders[i]["status"]
    order_status = {
        "order_id": i + 1,
        "status": "создан",
        "updated_at": order_created_at
    }
    order_statuses.append(order_status)
    if status == "создан":
        continue

    min_value = 1
    max_value = abs(order_updated_at - order_created_at).days
    if status == "отменен":
        if r.choice([True, False]):
            order_status = {
                "order_id": i + 1,
                "status": "отправлен",
                "updated_at": order_created_at + timedelta(days=fake.random_int(min=min_value, max=max_value))
            }
            order_statuses.append(order_status)
    elif status == "доставлен":
        order_status = {
            "order_id": i + 1,
            "status": "отправлен",
            "updated_at": order_created_at + timedelta(days=fake.random_int(min=min_value, max=max_value))
        }
        order_statuses.append(order_status)
    
    order_status = {
        "order_id": i + 1,
        "status": status,
        "updated_at": order_updated_at
    }
    order_statuses.append(order_status)

orders_couriers = []
pairs_order_courier = r.sample(list(itertools.product(range(1, len(couriers)), range(1, len(orders)))), records_num)
for pair in pairs_order_courier:
    min_value = 1
    max_value = abs(orders[pair[1]]["updated_at"] - orders[pair[1]]["created_at"]).days
    if max_value < min_value:
        assigned_at = orders[pair[1]]["created_at"]
    else:
        assigned_at = orders[pair[1]]["created_at"] + timedelta(days=fake.random_int(min=min_value, max=max_value))

    order_courier = {
        "courier_id": pair[0],
        "order_id": pair[1],
        "assigned_at": assigned_at
    }
    orders_couriers.append(order_courier)


with open("./data/users.csv", "w") as f:
    writer = csv.DictWriter(f, fieldnames=["user_type", "name", "phone", "email"])
    writer.writeheader()
    for user in users:
        writer.writerow(user)

with open("./data/couriers.csv", "w") as f:
    writer = csv.DictWriter(f, fieldnames=["name", "surname", "patronymic", "phone", "status"])
    writer.writeheader()
    for courier in couriers:
        writer.writerow(courier)

with open("./data/addresses.csv", "w") as f:
    writer = csv.DictWriter(f, fieldnames=["region", "city", "address_street", "address_type"])
    writer.writeheader()
    for address in addresses:
        writer.writerow(address)

with open("./data/orders.csv", "w") as f:
    writer = csv.DictWriter(f, fieldnames=["created_at", "updated_at", "status", "delivery_cost", "sender_user_id", "receiver_user_id", "pickup_address_id","delivery_address_id"])
    writer.writeheader()
    for order in orders:
        writer.writerow(order)

with open("./data/orderItem.csv", "w") as f:
    writer = csv.DictWriter(f, fieldnames=["order_id", "item_name", "quantity"])
    writer.writeheader()
    for order_item in order_items:
        writer.writerow(order_item)

with open("./data/orderCourier.csv", "w") as f:
    writer = csv.DictWriter(f, fieldnames=["courier_id", "order_id", "assigned_at"])
    writer.writeheader()
    for order_courier in orders_couriers:
        writer.writerow(order_courier)

with open("./data/orderStatus.csv", "w") as f:
    writer = csv.DictWriter(f, fieldnames=["order_id", "status", "updated_at"])
    writer.writeheader()
    for order_status in order_statuses:
        writer.writerow(order_status)




