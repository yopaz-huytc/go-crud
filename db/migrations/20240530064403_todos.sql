-- migrate:up
create table todos (
       id integer primary key AUTO_INCREMENT not null,
       name varchar(255) not null,
       description text default null,
       is_done tinyint(1) default 0,
       created_at timestamp not null default current_timestamp,
       updated_at timestamp not null default current_timestamp on update current_timestamp,
       deleted_at timestamp null default null
);

-- migrate:down
drop table todos;
