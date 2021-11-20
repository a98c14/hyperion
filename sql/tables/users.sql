create table users (
  id           bigserial primary key,
  name         varchar(50) not null,
  created_date timestamp default current_timestamp
);