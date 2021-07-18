-- +goose Up
CREATE TABLE IF NOT EXISTS events (
            ID TEXT PRIMARY KEY,
            Title TEXT,
            DateTime timestamptz,
            Duration BIGINT,
            Description TEXT,
            UserID TEXT,
            NoteBefore BIGINT
);

INSERT INTO events (ID, Title, DateTime, Duration, Description, UserID, NoteBefore)
VALUES
    ('bp', 'HBD', '2021-09-15 17:00:12+03', '1', 'birthday party', 'vva', '2'),
    ('chil', 'chilling', '2021-07-15 12:01:12+03', '5000000', 'vacation chil', 'vva', '6000000'),
    ('vac', 'vacation', '2021-07-14 12:00:12+03', '3', 'vacation trip', 'vva', '4');

-- +goose Down
drop table events;

