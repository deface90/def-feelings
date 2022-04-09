CREATE TABLE IF NOT EXISTS "user"
(
    id                     serial primary key,
    username               varchar(50)              not null default '',
    password               varchar(255)             not null default '',
    full_name              varchar(100)             not null default '',
    email                  varchar(100)             not null default '',
    tg_username            varchar(100)             not null default '',
    notification_type      smallint                 not null default 0,
    status                 smallint                 not null default 0,
    created                timestamp with time zone not null default now(),
    last_login             timestamp with time zone not null default timestamp 'epoch',
    settings               json,
    auth_token             varchar(255)             not null default '',
    last_notification      timestamp with time zone not null default timestamp 'epoch',
    notification_frequency int                      not null default 0
);
INSERT INTO "user"
VALUES (1, 'test', 'test', 'test testing', 'test@test.com', 'tg_test', 1, 4, now(), now(), '{}');
SELECT setval('user_id_seq', (SELECT max(id) FROM "user"));

CREATE TABLE IF NOT EXISTS "user_subscription"
(
    id                serial primary key,
    type              integer                  not null default 0,
    user_id           integer                  not null references "user" on delete cascade,
    notification_type smallint                 not null default 0,
    subscription_id   varchar(100)             not null default '',
    chat_id           int                      not null default 0,
    last_notification timestamp with time zone not null default timestamp 'epoch',
    error             varchar(100)             not null default '',
    status            smallint                 not null default 0
);
INSERT INTO "user_subscription"
VALUES (1, 1, 1, 0, '123', 456, now(), '', 1);
SELECT setval('user_subscription_id_seq', (SELECT max(id) FROM "user_subscription"));

CREATE TABLE IF NOT EXISTS "feeling"
(
    id    serial primary key,
    title varchar(255) not null unique
);
INSERT INTO feeling
VALUES (1, 'foo');
SELECT setval('feeling_id_seq', (SELECT max(id) FROM "feeling"));

CREATE TABLE IF NOT EXISTS "status"
(
    id      serial primary key,
    user_id integer                  references "user" on delete set null,
    message text                     not null default '',
    created timestamp with time zone not null default now()
);
insert into status
values (1, 1, 'test message', now());
SELECT setval('status_id_seq', (SELECT max(id) FROM "status"));

CREATE TABLE IF NOT EXISTS "status_feeling"
(
    status_id  integer not null references status on delete cascade,
    feeling_id integer not null references feeling on delete cascade,
    primary key (status_id, feeling_id)
);
insert into status_feeling
values (1, 1);
