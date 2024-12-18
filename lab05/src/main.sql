-- 1 извлечь json фрагмент из json документа
select 
       name || ' ' || surname as full_name,
       additional_data->>'experience' as experience_years,
       additional_data->'vehicle' as vehicle_info
from couriers_json_atr;

-- 2 извлечь значения конкретных узлов или атрибутов json документа
select 
       name || ' ' || surname as full_name,
       additional_data->>'experience' as experience_years,
       additional_data->'vehicle'->>'type' as vehicle_type,
       additional_data->'vehicle'->>'number' as vehicle_number
from couriers_json_atr;

-- 3 выполнить проверку существования узла или атрибута
insert into couriers_json_atr (name, surname, patronymic, phone, status, additional_data) values
('Алексей', 'Сидоров', 'Петрович', '123-123-1234', 'занят', null),
('Максим', 'Кузнецов', null, '321-321-4321', 'свободен', '{"experience": 2, "regions": ["Казань"], "vehicle": "null"}');

select *
from couriers_json_atr
where additional_data is not null;

select *
from couriers_json_atr
where additional_data is not null and additional_data->'vehicle' is not null;

-- 4 изменить json документ
update couriers_json_atr
set additional_data = '{"experience": 5, "regions": ["москва", "тула"], "vehicle": {"type": "мотоцикл", "number": "c789de77"}}'
where additional_data is null;

-- 5 разделить json документ на несколько строк по узлам
drop table if exists json_table;
create table if not exists json_table (
    data jsonb
);

\copy json_table(data) from '/home/void/student/bmstu/sem5/db/lab05/src/data/couriers.json';

select jsonb_array_elements(data)
from json_table;
