CREATE TABLE IF NOT EXISTS "user"
(
    id                serial primary key,
    username          varchar(50)              not null default '',
    password          varchar(255)             not null default '',
    full_name         varchar(100)             not null default '',
    email             varchar(100)             not null default '',
    tg_username       varchar(100)             not null default '',
    notification_type smallint                 not null default 0,
    status            smallint                 not null default 0,
    created           timestamp WITH TIME ZONE not null default now(),
    last_login        timestamp WITH TIME ZONE not null default timestamp 'epoch',
    settings          json
);
CREATE TABLE IF NOT EXISTS "feeling"
(
    id    serial primary key,
    title varchar(255) not null unique
);
CREATE TABLE IF NOT EXISTS "status"
(
    id      serial primary key,
    user_id integer references "user" on delete set null,
    message text    not null default ''
);
CREATE TABLE IF NOT EXISTS "status_feeling"
(
    status_id  integer not null references status on delete cascade,
    feeling_id integer not null references feeling on delete cascade,
    primary key (status_id, feeling_id)
);
