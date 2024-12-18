create or replace function get_orders_with_items(item_names text[])
returns table(order_id int) as $$
begin
    return query
    select distinct oi.order_id
    from orderitem oi
    where oi.item_name = any(item_names);
end;
$$ language plpgsql;
