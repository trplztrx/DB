create table if not exists table1 (
    id int not null,
    var1 text not null,
    valid_from_dttm timestamp without time zone default current_timestamp,
    valid_to_dttm timestamp without time zone default current_timestamp,
    primary key (id, vlaid_from_dttm)
);

create table if not exists table2 (
    id int not null,
    var2 text not null,
    valid_from_dttm timestamp without time zone default current_timestamp,
    valid_to_dttm timestamp without time zone default current_timestamp,
    primary key (id, vlaid_from_dttm)
);

insert into table1(id, var1, valid_from_dttm, valid_to_dttm)
values (1, 'a', '2018-09-01', '2018-09-15');

insert into table1(id, var1, valid_from_dttm, valid_to_dttm)
values (1, 'b', '2018-09-16', '5999-12-31');

insert into table2(id, var2, valid_from_dttm, valid_to_dttm)
values (1, 'a', '2018-09-01', '2018-09-18');

insert into table2(id, var2, valid_from_dttm, valid_to_dttm)
values (1, 'b', '2018-09-18', '5999-12-31');

select 
    t1.id,
    t1.var1,
    t2.var2,
    greatest(t1.valid_from_dttm, t2.valid_from_dttm) as valid_from_dttm,
    least(t1.valid_to_dttm, t2.valid_to_dttm) as valid_to_dttm
from table1 t1
join table2 t2 on
    t1.id=t2.id
    and t1.valid_from_dttm <= t2.valid_to_dttm
    and t1.valid_to_dttm >= t2.valid_from_dttm
order by id, valid_from_dttm