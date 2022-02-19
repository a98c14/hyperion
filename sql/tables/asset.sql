create table asset (
  id                serial primary key,
  unity_guid        varchar(100) not null,
  unity_internal_id bigint not null,
  name              varchar(200) not null,
  guid              uuid not null default uuid_generate_v4(),

  -- Type of the asset. Sprite, material, animation, prefab etc.
  type     integer not null,
  
  deleted_date timestamp,
  created_date timestamp default current_timestamp
);
