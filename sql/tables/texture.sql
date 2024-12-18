/*
  Stores the textures used in Unity. Referenced by sprites.
*/
create table texture (
  id           serial primary key,
  asset_id     integer not null,
  image_path   varchar(500) not null,
  created_date timestamp default current_timestamp,
  foreign key (asset_id) references asset(id)
);
