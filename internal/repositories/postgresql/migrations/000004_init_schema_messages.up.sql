CREATE TABLE IF NOT EXISTS "messages" (
             message_id SERIAL PRIMARY KEY,
             room_id varchar(255) REFERENCES rooms(room_id) not null ,
             sender_id integer REFERENCES users(id)not null,
             content TEXT NOT NULL,
             created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS "messages_statuses"(
            message_id INTEGER REFERENCES messages(message_id)not null,
            user_id INTEGER REFERENCES users(id)not null ,
            is_read BOOLEAN DEFAULT FALSE not null,
            read_at TIMESTAMP ,
            PRIMARY KEY (message_id, user_id)

);
CREATE INDEX ON messages(room_id, created_at);
CREATE INDEX ON messages_statuses(user_id, is_read);