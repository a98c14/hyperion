create table tag (
  id           serial primary key,
  name         varchar(200) not null,
  color        char(7) not null,
  deleted_date timestamp,
  created_date timestamp default current_timestamp
);
