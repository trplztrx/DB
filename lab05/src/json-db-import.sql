drop table if exists order_items_json;
drop table if exists json_table;

create temporary table if not exists order_items_json (
    id serial primary key,
    order_id int not null,
    item_name varchar(255) not null,
    quantity int not null,
    constraint chk_quantity check (quantity > 0)
);

create table if not exists json_table (
    data jsonb
);

\copy json_table(data) from '/home/void/student/bmstu/sem5/db/lab05/src/data/orderItem.json';

with json_elements AS (
    select jsonb_array_elements(data) as d
    from json_table
)
insert into order_items_json (id, order_id, item_name, quantity)
select
    (d::jsonb->>'id')::int as id,
    (d::jsonb->>'order_id')::int as order_id,
    d::jsonb->>'item_name' as item_name,
    (d::jsonb->>'quantity')::int as quantity
from json_elements;

select * from order_items_json;
