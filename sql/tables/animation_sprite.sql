/*
  Stores the textures used in Unity. Referenced by sprites.
*/
create table animation_sprite (
  id                 serial primary key,
  animation_id       integer not null,
  sprite_id          integer not null,
  created_date       timestamp default current_timestamp,
  foreign key (animation_id) references animation(id),
  foreign key (sprite_id) references sprite(id)
);
