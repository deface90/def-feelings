ALTER TABLE "user" ADD COLUMN last_notification timestamp WITH TIME ZONE NOT NULL DEFAULT now();
ALTER TABLE "user" ADD COLUMN notification_frequency int not null default 0;