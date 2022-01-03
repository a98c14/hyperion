-- public int Priority; 
--         public Sprite[] Sprites;            
--         public AnimationTransition TransitionType;
create table animation (
  id              serial primary key,
  name            varchar(200) not null,
  priority        integer not null,
  transition_type integer not null,
  created_date    timestamp default current_timestamp
);
