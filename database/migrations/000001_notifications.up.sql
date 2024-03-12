CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    identifier TEXT NOT NULL,
    subject TEXT NOT NULL DEFAULT '',
    message TEXT NOT NULL,
    type SMALLINT NOT NULL,
    attempts SMALLINT DEFAULT 0,
    created_at INT NOT NULL,
    sent_at INT DEFAULT NULL
);

CREATE INDEX notifications_type_attempts ON notifications (type, attempts);
