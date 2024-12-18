drop trigger if exists after_order_update on orders;
drop function if exists log_order_update;

drop view if exists order_view;
drop function if exists order_view_func;
drop trigger if exists order_view_trigger on order_view;