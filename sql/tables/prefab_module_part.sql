/*
Stores the actual gameplay values for prefabs

Example: 
Creature module has a field named OffensiveStats. OffensiveStats has a field named Health.
To create a prefab named `EnemyBoar` with health value 100. Below rows should exist in database.

    SampleBalance                  - Balance Version Id: 1
    EnemyBoar                      - Prefab Id: 1
    Creature                       - Prefab Part Id: 5
    Creature.OffensiveStats.Health - Prefab Part Id: 17

    PrefabPartValue:
        - Prefab Part Id: 17
        - Prefab Id: 1
        - Balance Version Id: 1
        - Array Index: 0
        - Value: 100
*/
create table prefab_module_part (
    id                 bigserial primary key,

    -- If referenced module value has array type this column 
    -- determines the index. For non array types has the value of 0
    array_index        integer not null, 

    -- Value for the referenced field. Either numeric or string value is used
    -- depending on the type.
    value_type         integer not null,
    value              json,

    prefab_id          integer not null,
    module_part_id     integer not null,
    balance_version_id integer not null,
    created_date       timestamp default current_timestamp,
    unique (prefab_id, module_part_id, array_index, balance_version_id),
    foreign key (prefab_id)          references prefab(id),
    foreign key (module_part_id)     references module_part(id),
    foreign key (balance_version_id) references balance_version(id)
);
