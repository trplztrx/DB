CREATE TABLE IF NOT EXISTS orders (
    id SERIAL, -- PK
    created_at TIMESTAMP,
    delivery_date TIMESTAMP,
    status VARCHAR(50), -- создан, отправлен, доставлен
    delivery_cost DECIMAL(10, 2),
    sender_user_id INT, -- FK к user
    receiver_user_id INT, -- FK к user
    pickup_address_id INT, -- FK к address
    delivery_address_id INT -- FK к address
);

CREATE TABLE IF NOT EXISTS users(
    id SERIAL, -- PK
    user_type VARCHAR(50), -- Физ / Юр лицо
    name VARCHAR(255), -- ФИО / Компания
    phone VARCHAR(50),
    email VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL, -- PK
    region VARCHAR(100),
    city VARCHAR(100),
    street VARCHAR(100),
    house VARCHAR(20),
    apartments VARCHAR(20),
    address_type VARCHAR(50) -- отправитель / получатель
);

CREATE TABLE IF NOT EXISTS couriers (
    id SERIAL, -- PK
    name VARCHAR(255),
    surname VARCHAR(255),
    patronymic VARCHAR(255),
    phone VARCHAR(50),
    status VARCHAR(50) -- свободен, в пути, занят
);

CREATE TABLE IF NOT EXISTS orderItem (
    id SERIAL, -- PK
    order_id INT, -- FK к order
    item_name VARCHAR(255),
    quantity INT
);

CREATE TABLE IF NOT EXISTS orderStatus (
    id SERIAL, -- PK
    order_id INT, -- FK к order
    status VARCHAR(50),
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orderCourier (
    courier_id INT, -- FK к таблице courier
    order_id INT, -- FK к таблице order
    assigned_at TIMESTAMP -- время назначения курьера на заказ
);