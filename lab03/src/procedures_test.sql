-- 1
call add_or_update_order_item(3, 'товар_А', 10);

-- 2 
call generate_and_save_temp_dates(10);
select * from temp_generated_dates;

-- 3
call assign_courier_to_order(402);

select *
from orders
where id = 402;

-- 4
call get_metadata('public');