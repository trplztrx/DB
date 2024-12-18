-- 1. Инструкция SELECT, использующая предикат сравнения.
-- Список отправителей из Москвы
select u.name, u.user_type, a.city
from users u
join orders o on o.sender_user_id = u.id
join addresses a on o.pickup_address_id = a.id
where a.city = 'Москва'
-- 2. Инструкция SELECT, использующая предикат BETWEEN.
-- Список заказов созданных в период с 2000 до 2010 года
select *
from Orders
where created_at between '2000-01-01' AND '2010-01-01';
-- 3. Инструкция SELECT, использующая предикат LIKE.
-- Список курьеров у которых Фамилия начинается на Иван
select *
from couriers 
where surname like 'Иван%';
-- 4. Инструкция SELECT, использующая предикат IN с вложенным подзапросом.
-- Все имена, номера телефонов, почты пользователей которым был доставлен заказ
select name, phone, email
from users
where id in (select receiver_user_id
            from orders
            where status = 'доставлен')
-- 5. Инструкция SELECT, использующая предикат EXISTS с вложенным
-- подзапросом.
-- Заказы которые были переданы курьерам
select *
from orders o 
where exists (select * from ordercourier oc where oc.order_id = o.id);
-- 6. Инструкция SELECT, использующая предикат сравнения с квантором.
-- Заказы, стоимость которых больше чем хотя бы у одного
select *
from orders o 
where o.delivery_cost > any (
						select o2.delivery_cost
						from orders o2
						where o2.id != o.id)
-- 7. Инструкция SELECT, использующая агрегатные функции в выражениях
-- столбцов.
-- сгруппированная по статусу для анлитики заказов
select
    status, 
    count(*) AS total_orders, -- Общее количество заказов по статусу
    avg(delivery_cost) AS avg_delivery_cost, -- Средняя стоимость доставки по статусу
    sum(delivery_cost) AS total_delivery_cost -- Общая стоимость доставки по статусу
from
    orders
group by
    status
order by
    total_orders desc;

-- 8. Инструкция SELECT, использующая скалярные подзапросы в выражениях
-- столбцов.
-- вывести список заказов с указанием имени отправителя и получателя
select 
    id as order_id,
    (select name from users where id = orders.sender_user_id) as sender_name,
    (select name from users where id = orders.receiver_user_id) as receiver_name,
    delivery_cost,
    status
from 
    orders
order by 
    created_at desc;

-- 9. Инструкция SELECT, использующая простое выражение CASE.
-- вывести список заказов, в котором статус заказа классифицируется как завершен, в работе, отменен
select 
    id as order_id,
    delivery_cost,
    case 
        when status = 'доставлен' then 'завершен'
        when status = 'создан' or status = 'отправлен' then 'в работе'
        when status = 'отменен' then 'отменен'
        else 'неизвестный статус'
    end as status_category
from 
    orders
order by 
    created_at desc;

-- 10. Инструкция SELECT, использующая поисковое выражение CASE.
-- разделить заказы на категории: "низкая стоимость", "средняя стоимость" и "высокая стоимость"
select 
    id as order_id,
    delivery_cost,
    case 
        when delivery_cost < 25000000 then 'низкая стоимость'
        when delivery_cost >= 25000000 and delivery_cost <= 75000000 then 'средняя стоимость'
        when delivery_cost > 75000000 then 'высокая стоимость'
        else 'стоимость не указана'
    end as cost_category
from 
    orders
order by 
    delivery_cost;

-- 11. Создание новой временной локальной таблицы из результирующего набора
-- данных инструкции SELECT.
-- создать временную таблицу с заказами, у которых доставка стоит более 50000000
create temporary table expensive_orders as
select 
    id as order_id,
    delivery_cost,
    status
from 
    orders
where 
    delivery_cost > 50000000;

-- 12. Инструкция SELECT, использующая вложенные коррелированные
-- подзапросы в качестве производных таблиц в предложении FROM.
-- получить список заказов с максимальной стоимостью доставки для каждого пользователя
select 
    users.id as user_id,
    users.name as user_name,
    max_delivery.max_cost as max_delivery_cost
from 
    users
join 
    (select 
         orders.sender_user_id,
         max(orders.delivery_cost) as max_cost
     from 
         orders
     group by 
         orders.sender_user_id
    ) as max_delivery on users.id = max_delivery.sender_user_id;

-- 13. Инструкция SELECT, использующая вложенные подзапросы с уровнем
-- вложенности 3.
-- найти заказы с максимальной стоимостью доставки среди пользователей, у которых общее количество товаров в заказах превышает 500
select 
    orders.id as order_id,
    orders.delivery_cost,
    orders.sender_user_id
from 
    orders
where 
    orders.delivery_cost = (
        select 
            max(delivery_cost)
        from 
            orders as inner_orders
        where 
            inner_orders.sender_user_id = (
                select 
                    sender_user_id
                from 
                    (
                        select 
                            sender_user_id,
                            sum(orderitem.quantity) as total_quantity
                        from 
                            orders
                        join 
                            orderitem on orders.id = orderitem.order_id
                        group by 
                            sender_user_id
                        having 
                            sum(orderitem.quantity) > 500
                    ) as high_quantity_users
                where 
                    high_quantity_users.sender_user_id = orders.sender_user_id
            )
    );

-- 14. Инструкция SELECT, консолидирующая данные с помощью предложения
-- GROUP BY, но без предложения HAVING.
-- сгруппированная по статусу для анлитики заказов
select
    status, 
    count(*) AS total_orders, -- Общее количество заказов по статусу
    avg(delivery_cost) AS avg_delivery_cost, -- Средняя стоимость доставки по статусу
    sum(delivery_cost) AS total_delivery_cost -- Общая стоимость доставки по статусу
from
    orders
group by
    status
order by
	total_orders desc;

-- 15. Инструкция SELECT, консолидирующая данные с помощью предложения
-- GROUP BY и предложения HAVING.
-- получить статусы заказов, где общее количество заказов превышает 250, и подсчитать их общее количество и суммарную стоимость доставки
select 
    status,
    count(*) as total_orders,
    sum(delivery_cost) as total_delivery_cost
from 
    orders
group by 
    status
having 
    count(*) > 250
order by 
    total_orders desc;

-- 16. Однострочная инструкция INSERT, выполняющая вставку в таблицу одной
-- строки значений.
insert into orders (created_at, updated_at, status, delivery_cost, sender_user_id, receiver_user_id, pickup_address_id, delivery_address_id) 
values ('2024-11-20 12:00:00', '2024-11-20 12:00:00', 'создан', 500.00, 1, 2, 10, 20);

-- 17. Многострочная инструкция INSERT, выполняющая вставку в таблицу
-- результирующего набора данных вложенного подзапроса.
-- добавляет запись в таблицу orderitem только для заказов, где товара с названием 'новый_товар' ещё нет.
insert into orderitem (order_id, item_name, quantity)
select 
    orders.id as order_id,
    'новый_товар' as item_name,
    1 as quantity
from 
    orders
where 
    orders.id not in (
        select 
            order_id
        from 
            orderitem
        where 
            item_name = 'новый_товар'
    );

--18. Простая инцструкция Update
update orderitem
set item_name = 'не_новый_товар'
where item_name = 'новый_товар';
-- 19. Инструкция UPDATE со скалярным подзапросом в предложении SET.
-- UPDATE Products
-- обновить стоимость доставки для заказов
update orders
set delivery_cost = (
    select min(delivery_cost)
    from orders
)
where status = 'создан';

-- 20. Простая инструкция DELETE.
delete from orderitem
where item_name = 'не_новый_товар';
-- 21. Инструкция DELETE с вложенным коррелированным подзапросом в
-- предложении WHERE.
-- все строки из таблицы orderitem, где количество товара превышает среднее для соответствующего заказа
delete from orderitem
where quantity > (
    select avg(quantity)
    from orderitem as sub
    where sub.order_id = orderitem.order_id
);

-- 22. Инструкция SELECT, использующая простое обобщенное табличное
-- выражение
-- все заказы, где общее количество товаров превышает 50
with order_summary as (
    select 
        order_id,
        sum(quantity) as total_items
    from 
        orderitem
    group by 
        order_id
)
select 
    orders.id as order_id,
    orders.status,
    order_summary.total_items
from 
    orders
join 
    order_summary on orders.id = order_summary.order_id
where 
    order_summary.total_items > 50;

-- 23. Инструкция SELECT, использующая рекурсивное обобщенное табличное
-- выражение.
-- список дат от текущей даты до 10 дней вперёд
with recursive delivery_dates as (
    -- начинаем с текущей даты
    select 
        current_date as delivery_date,
        1 as day_number
    union all
    -- рекурсивный случай: добавляем один день за раз
    select 
        (delivery_date + interval '1 day')::date as delivery_date,
        day_number + 1
    from 
        delivery_dates
    where 
        day_number < 10  -- генерируем даты на 10 дней вперёд
) 
select 
    delivery_date
from 
    delivery_dates;

-- 24. Оконные функции. Использование конструкций MIN/MAX/AVG OVER()
-- min/max/avg cатистика по стоимости доставки заказов
select 
    id as order_id,
    status,
    delivery_cost,
    min(delivery_cost) over (partition by status) as min_delivery_cost,
    max(delivery_cost) over (partition by status) as max_delivery_cost,
    avg(delivery_cost) over (partition by status) as avg_delivery_cost
from 
    orders
order by 
    status, delivery_cost;

-- 25. Оконные фнкции для устранения дублей
-- Придумать запрос, в результате которого в данных появляются полные дубли.
-- Устранить дублирующиеся строки с использованием функции ROW_NUMBER().
select 
    id as order_id,
    status,
    delivery_cost
from 
    orders
union all
select 
    id as order_id,
    status,
    delivery_cost
from 
    orders;

with numbered_rows as (
    select 
        id as order_id,
        status,
        delivery_cost,
        row_number() over (partition by id, status, delivery_cost order by id) as row_num
    from 
        orders
)
select 
    order_id,
    status,
    delivery_cost
from 
    numbered_rows
where 
    row_num = 1;


