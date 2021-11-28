create table component_part (
  id           serial primary key,
  name         varchar(200) not null,
  created_date timestamp default current_timestamp,
  foreign key (parent_id) references component(id)
);
