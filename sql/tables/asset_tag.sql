create table asset_tag (
    asset_id     integer not null,
    tag_id       integer not null,
    deleted_date timestamp,
    created_date timestamp default current_timestamp,
    unique (asset_id, tag_id),
    foreign key (tag_id) references tag(id),
    foreign key (asset_id) references asset(id)
);
