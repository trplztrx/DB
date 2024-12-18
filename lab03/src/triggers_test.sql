-- 1
update orders
set status = 'отправлен', updated_at = current_timestamp
where id = 1;

-- 2
insert into order_view (status, delivery_cost, item_name, quantity)
values ('создан', 5000000000.00, 'Товар А', 10);