create table prefab (
  id           serial primary key,
  asset_id     integer not null,
  parent_id    integer,
  transform    jsonb not null,
  renderer     jsonb not null,
  colliders    jsonb not null,
  created_date timestamp default current_timestamp,
  foreign key (parent_id) references prefab(id),
  foreign key (asset_id) references asset(id)
);
