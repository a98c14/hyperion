create table prefab (
  id           serial primary key,
  name         varchar(200) not null,
  parent_id    integer,
  created_date timestamp default current_timestamp,
  foreign key (parent_id) references prefab(id)
);
