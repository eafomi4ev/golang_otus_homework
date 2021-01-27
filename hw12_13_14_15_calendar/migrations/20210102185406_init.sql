-- +goose Up
create table if not exists events
(
    id          varchar(15) primary key,
    title       varchar(150) not null default '',
    event_date  timestamptz  not null default now(),
    duration    integer      not null default 0,
    description text         not null default '',
    user_id     varchar(36)  not null,
    remind_in   integer               default 0
);

-- +goose Down
drop table events;
