CREATE TABLE IF NOT EXISTS rooms
(
    room_id varchar(255) not null primary key
);

CREATE INDEX on rooms (room_id);

