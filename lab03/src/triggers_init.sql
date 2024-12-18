-- триггер after для записи в orderStatus обновления поля status отношения orders
create or replace function log_order_update()
returns trigger as $$
begin
    if new.status is distinct from old.status then
        insert into orderStatus (order_id, status, updated_at)
        values (new.id, new.status, new.updated_at)
    end if;

    return new;
end;
$$ language plpgsql;

create trigger after_order_update
after update on orders
for each row
execute function log_order_update();

-- DML триггер intead of для упрощенной вставки данных через пердставление
create view order_view as
select 
    o.id as order_id,
    o.status,
    o.delivery_cost,
    oi.item_name,
    oi.quantity
from orders o
left join orderitem oi on o.id = oi.order_id;

create or replace function order_view_func()
returns trigger as $$
begin
    -- вставка данных в таблицу orders с заполнением обязательных полей
    insert into orders (created_at, updated_at, status, delivery_cost)
    values (
        now(),
        now(),
        new.status,
        new.delivery_cost
    )
    returning id into new.order_id;

    -- если указаны данные для товаров, вставляем их в таблицу orderitem
    if new.item_name is not null and new.quantity is not null then
        insert into orderitem (order_id, item_name, quantity)
        values (new.order_id, new.item_name, new.quantity);
    end if;

    return new;
end;
$$ language plpgsql;


create trigger order_view_trigger
instead of insert
on order_view
for each row
execute function order_view_func();

