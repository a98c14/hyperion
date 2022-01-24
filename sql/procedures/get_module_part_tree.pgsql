with recursive module_part_recursive as (
    select 
        id, 
        name, 
        value_type, 
        parent_id, 
        case when parent_id is not null then name else null end as parent_name,
        is_array,
        tooltip
    from module_part
    where id=33 and parent_id is null and deleted_date is null
    union select 
        c.id, 
        c.name, 
        c.value_type, 
        c.parent_id, 
        (select name from module_part where id=cp.id),
        c.is_array,
        c.tooltip
    from module_part c inner join module_part_recursive cp on cp.id=c.parent_id
    where deleted_date is null)
select * from module_part_recursive;