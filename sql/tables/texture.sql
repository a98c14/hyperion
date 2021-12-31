/*
  Stores the textures used in Unity. Referenced by sprites.
*/
create table texture (
  id           serial primary key,
  image_path   varchar(max) not null,
  unity_guid   varchar(max) not null,
  unity_name   varchar(200) not null,
  created_date timestamp default current_timestamp
);
