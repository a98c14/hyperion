create table prefab (
  id           serial primary key,
  asset_id     integer not null,
  parent_id    integer,
  created_date timestamp default current_timestamp,
  foreign key (parent_id) references prefab(id),
  foreign key (asset_id) references asset(id)
);
