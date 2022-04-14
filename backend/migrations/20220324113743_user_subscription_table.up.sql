CREATE TABLE IF NOT EXISTS "user_subscription"
(
    id                serial primary key,
    type              integer                  not null default 0,
    user_id           integer                  not null references "user" on delete cascade,
    notification_type smallint                 not null default 0,
    subscription_id   varchar(100)             not null default '',
    chat_id           integer                  not null default 0,
    last_notification timestamp WITH TIME ZONE not null default timestamp 'epoch',
    error             varchar(100)             not null default '',
    status            smallint                 not null default 0
);