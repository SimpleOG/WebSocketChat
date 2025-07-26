-- name: CreateUser :one
INSERT INTO users (username,
                   password,
                   email,
                   date_of_birth)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByUsernameOrEmail :one
SELECT *
FROM users
WHERE username = $1 or email=$2;

-- name: GetUserForLogin :one
SELECT *
FROM users
WHERE (username = $1 or email = $2)
  and password = $3;

-- name: GetUsersById :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUsersWithLessDate :many
SELECT *
FROM users
WHERE date_of_birth < $1;