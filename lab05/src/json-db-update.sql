
drop table if exists couriers_json_atr;

create table if not exists couriers_json_atr (
    id serial primary key,
    name varchar(255) not null,
    surname varchar(255) not null,
    patronymic varchar(255),
    phone varchar(50) not null unique,
    status varchar(50) not null check (status in ('свободен', 'в пути', 'занят')),
    additional_data jsonb
);

insert into couriers_json_atr (name, surname, patronymic, phone, status, additional_data) values
('Иван', 'Иванов', 'Иванович', '123-456-7890', 'свободен', '{"experience": 5, "regions": ["Москва", "Санкт-Петербург"], "vehicle": {"type": "автомобиль", "number": "A123BC77"}}'),
('Петр', 'Петров', NULL, '987-654-3210', 'в пути', '{"experience": 3, "regions": ["Казань", "Нижний Новгород"], "vehicle": {"type": "велосипед"}}');

update couriers_json_atr
set additional_data = '{"experience": 7, "regions": ["Москва", "Санкт-Петербург", "Казань"], "vehicle": {"type": "грузовик", "number": "B456CD77"}}'
where phone = '123-456-7890';

select * from couriers_json_atr;
