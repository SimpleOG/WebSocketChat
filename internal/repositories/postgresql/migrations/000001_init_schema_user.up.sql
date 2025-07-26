CREATE TABLE if not exists "users"
(
    "id"            SERIAL              NOT NULL UNIQUE primary key,
    "username"      VARCHAR(255) UNIQUE NOT NULL,
    "email"         VARCHAR(255) UNIQUE,
    "date_of_birth" VARCHAR(255)        NOT NULL,
    "password"      VARCHAR(255)        NOT NULL,
    "created_at"    date                NOT NULL DEFAULT (now())
);



