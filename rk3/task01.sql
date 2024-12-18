-- РК 3 Шишков Константин ИУ7-55Б

create table satellite (
    "ID Спутника" serial primary key,
    "Название спутника" varchar(50) not null,
    "Дата производства" date not null,
    "Страна" varchar(50) not null
);

insert into satellite ("ID Спутника", "Название спутника", "Дата производства", "Страна") values
(1, 'SIT-2086', '2050-01-01', 'Россия'),
(2, 'Шицзянь 16-02', '2049-12-01', 'Китай');

-- Сделал PK отдельным атрибутом, т к посчитал усложнением составной PK из 3 полей. Это осмысленное действие.
create table flight (
	id serial primary key,
    "ID Спутника" int references satellite("ID Спутника"),
    "Дата запуска" date not null,
    "Время запуска" time not null,
    "День недели" varchar(20) not null,
    "Тип" int not null check ("Тип" in (0, 1))
);

insert into flight ("ID Спутника", "Дата запуска", "Время запуска", "День недели", "Тип") values
(1, '2050-05-11', '09:00', 'Среда', 1),
(1, '2051-06-14', '23:05', 'Среда', 0),
(1, '2051-10-10', '23:50', 'Вторник', 1),
(2, '2050-05-11', '15:15', 'Среда', 1),
(1, '2052-01-01', '12:15', 'Понедельник', 0);

drop table if exists satellite;
drop table if exists flight;

-- Запрос 1
-- Поиск спутников из таблицы satellite, которые не имеют полетов в таблице flight
select s."ID Спутника", s."Название спутника"
from satellite s
where not exists (
    select 1
    from flight f
    where f."ID Спутника" = s."ID Спутника"
);

-- Запрос 2
-- Соединение таблиц satellite и flight по условию совпадения месяца и года клонок "Дата производства" и "Дата запуска"
select f."ID Спутника", s."Название спутника", f."Дата запуска", s."Дата производства"
from flight f
join satellite s on date_trunc('month', f."Дата запуска"::timestamp) = date_trunc('month', s."Дата производства"::timestamp);
