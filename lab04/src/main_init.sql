create extension plpython3u;

-- Скалярная функция
-- Конкатенация ФИО курьера
create or replace function get_courier_full_name(courier_id int)
returns varchar
as $$
import plpy
query = """
    select surname, name, patronymic
    from couriers
    where id = %s
"""
result = plpy.execute(query % (courier_id,))
if result
    surname = result[0]['surname']
    name = result[0]['name']
    patronymic = result[0]['patronymic']
    return f"{surname} {name} {patronymic}"
else:
    return 'Курьер не найден'
$$ language plpython3u;

select get_courier_full_name(1);

select id, get_courier_full_name(id) as full_name
from couriers;

drop function if exists get_courier_full_name;

-- Агрегатная функция
-- Количество заказов с заданным статусом
create or replace function count_orders_by_status_sfunc(state int, value text, target_status text)
returns int
as $$
if value == target_status:
    return state + 1
return state
$$ language plpython3u;

create aggregate count_orders_by_status(text, text) (
    sfunc = count_orders_by_status_sfunc,
    stype = int,
    initcond = '0'
);

select count_orders_by_status(status, 'создан') as созданные_заказы,
       count_orders_by_status(status, 'доставлен') as доставленные_заказы
from orders;

drop aggregate if exists count_orders_by_status(text, text);
drop function if exists count_orders_by_status_sfunc(int, text, text);

-- Табличная функция
-- Возврат заказов по статусу
create or replace function get_orders_by_status(target_status text)
returns table(
    id int,
    created_at timestamp,
    updated_at timestamp,
    delivery_cost numeric,
    sender_user_id int,
    receiver_user_id int
)
as $$
query = """
    select id, created_at, updated_at, delivery_cost, sender_user_id, receiver_user_id
    from orders
    where status = %s
"""
result = plpy.execute(query % ("'" + target_status + "'"))

return result
$$ language plpython3u;

select * 
from get_orders_by_status('создан');

drop function if exists get_orders_by_status;

-- Хранимая процедура
-- Процедура будет добавлять новые товары в заказ, при этом увеличивать количество, если товар с таким именем уже существует для указанного заказа
create or replace procedure add_or_update_order_item(order_id_input int, item_name_input text, quantity_input int)
as $$
if quantity_input <= 0:
    raise Exception('Некорректный ввод')

query_exists = """
    select 1
    from orderitem
    where order_id = %s and item_name = %s
"""
result = plpy.execute(query_exists % (order_id_input, "'" + item_name_input + "'"))

if result:
    query_update = """
        update orderitem
        set quantity = quantity + %s
        where order_id = %s and item_name = %s
    """
    plpy.execute(query_update % (quantity_input, order_id_input, "'" + item_name_input + "'"))
else:
    query_insert = """
        insert into orderitem (order_id, item_name, quantity)
        values (%s, %s, %s)
    """
    plpy.execute(query_insert % (order_id_input, "'" + item_name_input + "'", quantity_input))
$$ language plpython3u;

call add_or_update_order_item(1, 'Товар_Б', 7);

call add_or_update_order_item(1, 'Товар_Б', 7);

select * from orderitem where order_id = 1;

drop function if exists add_or_update_order_item;

-- Триггер
-- Триггер after для записи в orderStatus обновления поля status отношения orders
create or replace function log_order_update()
returns trigger
as $$
if TD["new"]["status"] != TD["old"]["status"]:
    query = """
        insert into orderStatus (order_id, status, updated_at)
        values (%s, %s, %s)
    """
    plpy.execute(query % (TD["new"]["id"], "'" + TD["new"]["status"] + "'", "'" + str(TD["new"]["updated_at"]) + "'"))

return None
$$ language plpython3u;


create trigger after_order_update
after update on orders
for each row
execute function log_order_update();

update orders
set status = 'отменен', updated_at = current_timestamp
where id = 1;

select *
from orderstatus
where order_id = 1;

drop trigger if exists after_order_update on orders;
drop function if exists log_order_update;

-- Пользовательский тип данных
create type order_summary as (
    order_id int,
    status text,
    delivery_cost numeric
);

create or replace function get_order_summary(order_id_input int)
returns order_summary
as $$
query = """
    select id, status, delivery_cost
    from orders
    where id = %s
"""
result = plpy.execute(query % order_id_input)
if result:
    return (result[0]["id"], result[0]["status"], result[0]["delivery_cost"])
$$ language plpython3u;

select * from get_order_summary(370);

drop function if exists get_order_summary;
drop type if exists order_summary




