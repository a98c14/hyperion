/*
  Stores the textures used in Unity. Referenced by sprites.
*/
create table sprite (
  id                 serial primary key,
  texture_id         integer not null,
  unity_sprite_id    varchar(max) not null,
  unity_internal_id  varchar(max) not null,
  unity_name         varchar(200) not null,
  unity_pivot        json not null,
  unity_rect         json not null,
  unity_border       json not null,
  unity_alignment    integer not null,
  created_date       timestamp default current_timestamp,
  foreign key (texture_id) references texture(id)
);
