/*
    Entity Component Table 
    Type = { basic: 0, buffer: 1 }
*/
create table component (
  id           serial primary key,
  name         varchar(200) not null,
  type         integer not null,
  parent_id    integer,
  is_hidden    boolean,
  created_date timestamp default current_timestamp,
  foreign key (parent_id) references component(id)
);
