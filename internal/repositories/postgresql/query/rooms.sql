-- name: CreateRoom :exec
INSERT INTO rooms (room_id) values ($1);

