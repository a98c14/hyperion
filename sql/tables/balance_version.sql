create table balance_version (
  id           bigserial primary key,
  name         varchar(50) not null,
  deleted_date timestamp,
  created_date timestamp default current_timestamp
);