-- public int Priority; 
--         public Sprite[] Sprites;            
--         public AnimationTransition TransitionType;
create table animation (
  id              serial primary key,
  asset_id        integer not null,
  priority        integer not null,
  transition_type integer not null,
  created_date    timestamp default current_timestamp,
  foreign key (asset_id) references asset(id)
);
