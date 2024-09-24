import csv
import random as r
from faker import Faker
import re
import itertools
from datetime import timedelta, datetime

fake = Faker("ru_RU")

recordsNum = 1000
orderStatus = ["создан", "отправлен", "доставлен", "отменен"]
userType = ["физик", "юрик"]
addressType = ["отправитель", "получатель"]
courierType = ["свободен", "в пути", "занят"]

# print(fake.name())
ids_users = []
with open("./data/users.csv", "w") as f:
    writer = csv.writer(f, delimiter = ',', lineterminator='\n')

    writer.writerow(("id", "user_type", "name", "phone", "email"))
    for i in range (recordsNum):
        typeUser = r.choice(userType)
        if typeUser == "физик":
            name = fake.name()
        else:
            name = fake.company()    
        phone = fake.unique.phone_number()
        email = fake.unique.email()

        ids_users.append(i + 1)

        writer.writerow((i + 1, typeUser, name, phone, email))

ids_courier = []
with open("./data/couriers.csv", "w") as f:
    writer = csv.writer(f, delimiter=',', lineterminator='\n')
    writer.writerow(("id", "name", "surname", "patronymic", "phone", "status"))

    for i in range(recordsNum):
        name = fake.name().split()
        phone = fake.unique.phone_number()
        status = r.choice(courierType)

        ids_courier.append(i + 1)

        writer.writerow((i + 1, name[1], name[0], name[2], phone, status))

ids_pickup_address = []
ids_delivery_address = []
with open("./data/addresses.csv", "w") as f:
    writer = csv.writer(f, delimiter=',', lineterminator='\n')
    writer.writerow(("id", "region", "city", "street", "house", "apartments", "address_type"))

    pattern = r'(.+),\sд\.\s(\d+(?:/\d+)?)(?:\sк\.\s(\d+(?:/\d+)?))?(?:\sстр\.\s(\d+(?:/\d+)?))?'
    for i in range(recordsNum):
        address = fake.street_address()
        typeAddress = r.choice(addressType)
        region = fake.region()
        city = fake.city()
        match = re.match(pattern, address)

        if match:
            street = match.group(1)
            house_number = match.group(2)
            apartment_number = match.group(3)
        else:
            street = address
            house_number = None
            apartment_number = None

        if typeAddress == 'отправитель':
            ids_pickup_address.append(i + 1)
        else:
            ids_delivery_address.append(i + 1)

        writer.writerow((i + 1, region, city, street, house_number, apartment_number, typeAddress))

orders = []
with open("./data/orders.csv", "w") as f:
    writer = csv.writer(f, delimiter=',', lineterminator='\n')
    writer.writerow(("id", "created_at", "delivery_date", "status", "delivery_cost", "sender_user_id", "receiver_user_id", "pickup_address_id", "delivery_address_id"))

    for i in range(recordsNum):
        created_at = fake.date_time()
        status = r.choice(orderStatus)
        if status == "доставлен":
            delivery_date = created_at + timedelta(days=fake.random_int(min=1, max=30))
        else:
            delivery_date = None
        delivery_cost = fake.pydecimal(left_digits=8, right_digits=2, positive=True)

        pickup_address_id = r.choice(ids_pickup_address)
        delivery_address_id = r.choice(ids_delivery_address)

        sender_user_id, receiver_user_id = r.sample(ids_users, 2)

        orders.append((i + 1, created_at, delivery_date, status, delivery_cost, sender_user_id, receiver_user_id, pickup_address_id, delivery_address_id))

        writer.writerow((i + 1, created_at, delivery_date, status, delivery_cost, sender_user_id, receiver_user_id, pickup_address_id, delivery_address_id))

with open("./data/orderItem.csv", "w") as f:
    writer = csv.writer(f, delimiter=',', lineterminator='\n')
    writer.writerow(("id", "order_id", "item_name", "quantity"))

    r.shuffle(orders)
    id = 0
    for order in orders:
        order_id = order[0]
        for j in range(r.randint(1, 10)):
            id += 1
            item_name = fake.word()
            quantity = r.randint(1, 100)
            writer.writerow((id, order_id, item_name, quantity))

with open("./data/orderStatus.csv", "w") as f:
    writer = csv.writer(f, delimiter=',', lineterminator='\n')
    writer.writerow(("id", "order_id", "status", "updated_at"))

    id = 0
    for order in orders:
        id += 1
        order_id = order[0]
        updated_at = order[1]
        writer.writerow((id, order_id, orderStatus[0], updated_at))

        if order[3] == "отменен":
            if r.choice([True, False]):
                id += 1
                updated_at += timedelta(days=fake.random_int(min=1, max=7))
                writer.writerow((id, order_id, orderStatus[1], updated_at))
            
            id += 1
            updated_at += timedelta(days=fake.random_int(min=1, max=22))
            writer.writerow((id, order_id, orderStatus[3], updated_at))
        elif order[2] == None:
            id += 1
            updated_at += timedelta(days=fake.random_int(min=1, max=7))
            writer.writerow((id, order_id, orderStatus[1], updated_at))
        else:
            id += 1
            updated_at += timedelta(days=fake.random_int(min=1, max=(abs(updated_at - order[2]).days)))
            writer.writerow((id, order_id, orderStatus[1], updated_at))

            id += 1
            writer.writerow((id, order_id, orderStatus[2], order[2]))

with open("./data/orderCourier.csv", "w") as f:
    writer = csv.writer(f, delimiter=',', lineterminator='\n')
    writer.writerow(("courier_id", "order_id", "assigned_at"))

    ids_order = [order[0] for order in orders]
    ids_order_courier = r.sample(list(itertools.product(ids_order, ids_courier)), recordsNum)
    for pair in ids_order_courier:
        for item in orders:
            if item[0] == pair[1]:
                order = item
                break
        if order[1] == None:
            assigned_at = order[1] + timedelta(days = fake.random_int(min=1, max=(abs(order[2] - order[1]).days)-1))
        else:
            assigned_at = order[1] + timedelta(days = fake.random_int(min=1, max=30))
        writer.writerow((pair[0], pair[1], assigned_at))







