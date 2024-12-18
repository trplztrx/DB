-- скалярная функция, которая вычисляет общую стоимость доставки для заказов, сделанных конкретным отправителем
create or replace function get_total_delivery_cost(sender_id_input int)
returns numeric as $$
declare
    total_cost numeric;
begin
    select sum(delivery_cost)
    into total_cost
    from orders
    where sender_user_id = sender_id_input;

    if total_cost is null then
        return 0;
    end if;

    return total_cost;
end;
$$ language plpgsql;

-- подставляемая табличная функция, которая возвращает список заказо, сделанных конкретным отправителем
create or replace function get_orders_by_sender(sender_id_input int)
returns table(order_id int, delivery_cost numeric, status text) as $$
begin
    return query
    select
        orders.id as order_id,
        orders.delivery_cost,
        orders.status::text as status
    from
        orders
    where
        orders.sender_user_id = sender_id_input;
end;
$$ language plpgsql;

-- многооператорная табличная функция, которая возвращает общее кол-во заказов отправителя, суммарную стоимость доставки, детали каждого заказа
create or replace function get_sender_order_details(sender_id_input int)
returns table(order_id int, delivery_cost numeric, status text, updated_at timestamp, total_orders int, total_delivery_cost numeric) as $$
declare
    total_count int;
    total_cost numeric;
begin
    select count(*) into total_count
    from orders
    where sender_user_id = sender_id_input;

    select sum(orders.delivery_cost) into total_cost
    from orders
    where sender_user_id = sender_id_input;

    return query
    select
        orders.id as order_id,
        orders.delivery_cost,
        orders.status::text as status,
        orders.updated_at,
        total_count as total_orders,
        coalesce(total_cost, 0) as total_delivery_cost
    from
        orders
    where
        orders.sender_user_id = sender_id_input;

end;
$$ language plpgsql;

-- функция с рекурсивной ОТВ, возвращающая список дат от текущей даты до 10 дней вперёд
create or replace function generate_dates(days_ahead int)
returns table(delivery_date date) as $$
begin
    return query
    with recursive delivery_dates as (
        -- начинаем с текущей даты
        select 
            current_date as delivery_date,
            1 as day_number
        union all
        -- добавляем один день за раз
        select 
            (delivery_dates.delivery_date + interval '1 day')::date as delivery_date,
            delivery_dates.day_number + 1
        from 
            delivery_dates
        where 
            delivery_dates.day_number < days_ahead
    )
    select 
        delivery_dates.delivery_date
    from 
        delivery_dates;
end;
$$ language plpgsql;


