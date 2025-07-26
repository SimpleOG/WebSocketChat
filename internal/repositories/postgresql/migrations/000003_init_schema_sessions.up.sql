CREATE TABLE "sessions"
(
    id varchar(255) primary key not null,
    user_id integer not null,
    refresh_token varchar(512) not null ,
    is_revoked bool not null default false,
    created_at date default (now()),
    expires_at date

)