create table prefab (
  id           serial primary key,
  name         varchar(200) not null,
  created_date timestamp default current_timestamp
);
