\copy users(user_type, name, phone, email) from '/home/void/student/sem5/db/lab01/src/data/users.csv' delimiter ',' csv header;
\copy couriers(name,surname,patronymic,phone,status) from '/home/void/student/sem5/db/lab01/src/data/couriers.csv' delimiter ',' csv header;
\copy addresses(region,city,address_street,address_type) from '/home/void/student/sem5/db/lab01/src/data/addresses.csv' delimiter ',' csv header;
\copy orders(created_at,updated_at,status,delivery_cost,sender_user_id,receiver_user_id,pickup_address_id,delivery_address_id) from '/home/void/student/sem5/db/lab01/src/data/orders.csv' delimiter ',' csv header;
\copy orderCourier(courier_id,order_id,assigned_at) from '/home/void/student/sem5/db/lab01/src/data/orderCourier.csv' delimiter ',' csv header;
\copy orderItem(order_id,item_name,quantity) from '/home/void/student/sem5/db/lab01/src/data/orderItem.csv' delimiter ',' csv header;
\copy orderStatus(order_id,status,updated_at) from '/home/void/student/sem5/db/lab01/src/data/orderStatus.csv' delimiter ',' csv header;