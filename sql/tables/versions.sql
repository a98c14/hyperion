create table version (
  id           bigserial primary key,
  name         varchar(50) not null,
  content      jsonb not null,
  deleted_date timestamp,
  created_date timestamp default current_timestamp
);