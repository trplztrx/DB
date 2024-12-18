\copy (select array_to_json(array_agg(row_to_json(o))) from orders o) to '/home/void/student/bmstu/sem5/db/lab05/src/data/orders.json';
\copy (select array_to_json(array_agg(row_to_json(o))) from addresses o) to '/home/void/student/bmstu/sem5/db/lab05/src/data/addresses.json';
\copy (select array_to_json(array_agg(row_to_json(o))) from couriers o) to '/home/void/student/bmstu/sem5/db/lab05/src/data/couriers.json';
\copy (select array_to_json(array_agg(row_to_json(o))) from orderItem o) to '/home/void/student/bmstu/sem5/db/lab05/src/data/orderItem.json';
\copy (select array_to_json(array_agg(row_to_json(o))) from orderStatus o) to '/home/void/student/bmstu/sem5/db/lab05/src/data/orderStatus.json';
\copy (select array_to_json(array_agg(row_to_json(o))) from orderCourier o) to '/home/void/student/bmstu/sem5/db/lab05/src/data/orderCourier.json';