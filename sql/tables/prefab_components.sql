/*
*/
create table prefab_component (
    id           bigserial primary key,
    prefab_id    integer not null,
    component_id integer not null,
    buffer_index integer not null,
    created_date timestamp default current_timestamp,
    unique (prefab_id, component_id, buffer_index),
    foreign key (prefab_id)    references prefab(id),
    foreign key (component_id) references component(id)
);

create table prefab_component_revision (
    id                  bigserial primary key,
    prefab_component_id bigint not null,
    value               numeric,
    created_date        timestamp default current_timestamp,
    foreign key (prefab_component_id) references prefab_component(id)
);

/*
select * from prefab_component pc
inner join prefab_component_revision pcr on pcr.prefab_component_id=pc.id
*/