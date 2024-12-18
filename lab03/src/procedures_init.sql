-- процедура будет добавлять новые товары в заказ, при этом увеличивать количество, если товар с таким именем уже существует для указанного заказа
create or replace procedure add_or_update_order_item(order_id_input int, item_name_input text, quantity_input int)
as $$
begin
    if quantity_input <= 0 then
        raise exception 'Некорректный ввод';
    end if;

    if exists (
        select 1
        from orderitem
        where order_id = order_id_input and item_name = item_name_input
    ) then
        update orderitem
        set quantity = quantity + quantity_input
        where order_id = order_id_input and item_name = item_name_input;
    else
        insert into orderitem (order_id, item_name, quantity)
        values (order_id_input, item_name_input, quantity_input);
    end if;
end;
$$ language plpgsql ;

-- рекурсивная процедура генерирует список дат от текущей даты до указанного количества дней вперёд, а затем сохраняет их в таблицу generated_dates
create or replace procedure generate_and_save_temp_dates(days_ahead int)
language plpgsql as $$
begin
    create temporary table temp_generated_dates (
        id serial primary key,
        generated_date date not null
    );

    with recursive delivery_dates as (
        select 
            current_date as delivery_date,
            1 as day_number
        union all
        select 
            (delivery_date + interval '1 day')::date as delivery_date,
            day_number + 1
        from 
            delivery_dates
        where 
            day_number < days_ahead
    )
    insert into temp_generated_dates (generated_date)
    select 
        delivery_date
    from 
        delivery_dates;
end;
$$;

-- процедура с курсором, которая соединяет по id заказа курьера со статусом свободен
create or replace procedure assign_courier_to_order_with_cursor(order_id_input int)
as $$
declare
    -- Курсор для свободных курьеров
    courier_cursor cursor for
        select id, name
        from couriers
        where status = 'свободен';

    -- Переменные для хранения данных текущего курьера
    current_courier_id int;
    current_courier_name text;
begin
    if not exists (
        select 1
        from orders
        where id = order_id_input and status = 'создан'
    ) then
        raise notice 'Order ID % не валидный или не в статусе "создан".', order_id_input;
        return;
    end if;

    open courier_cursor;

    fetch courier_cursor into current_courier_id, current_courier_name;

    if not found then
        raise notice 'Нет свободных курьеров для назначения.';
        close courier_cursor;
        return;
    end if;

    insert into ordercourier (courier_id, order_id, assigned_at)
    values (current_courier_id, order_id_input, current_timestamp);

    update couriers
    set status = 'занят'
    where id = current_courier_id;

    update orders
    set status = 'отправлен',
        updated_at = current_timestamp
    where id = order_id_input;

    raise notice 'Courier ID %, name % назначен на заказ Order ID %', 
                 current_courier_id, current_courier_name, order_id_input;

    close courier_cursor;
end;
$$ language plpgsql;

-- Процедура для вывода информации об аргументах всех отношений схемы
create or replace procedure get_metadata(schema_name_input text)
as $$
declare
    -- Переменные для хранения данных о текущей строке
    current_table_name text;
    current_column_name text;
    current_data_type text;
    current_is_nullable text;
begin
    for current_table_name, current_column_name, current_data_type, current_is_nullable in
        select
            table_name,
            column_name,
            data_type,
            is_nullable
        from
            information_schema.columns
        where
            table_schema = schema_name_input
        order by
            table_name, ordinal_position
    loop
        raise notice 'Table: %, Column: %, Type: %, Nullable: %',
                     current_table_name, current_column_name, current_data_type, current_is_nullable;
    end loop;
end;
$$ language plpgsql;
