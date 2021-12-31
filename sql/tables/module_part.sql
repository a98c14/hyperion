/*
  Stores the prefab module structures. If a row has no parent id
  it means it is a root node. If a row has no id that references it.
  It means it is a leaf node. Only leaf nodes can have module values.
*/
create table module_part (
  id                 serial primary key,
  name               varchar(200) not null,
  value_type         integer not null,
  game_version_id    integer not null,  
  parent_id          integer,
  deleted_date       timestamp,
  created_date       timestamp default current_timestamp,
  foreign key (parent_id) references module_part(id)
);
