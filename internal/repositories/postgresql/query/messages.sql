-- name: CreateMessage :exec
INSERT INTO messages (room_id, sender_id, content)
VALUES ($1, $2, $3)
;

-- name: GetMessage :one
SELECT * FROM messages
WHERE message_id = $1;

-- name: ListRoomMessages :many
SELECT * FROM messages
WHERE room_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteMessage :exec
DELETE FROM messages
WHERE message_id = $1;

-- name: GetMessageWithStatus :one
SELECT
    m.*,
    ms.is_read,
    ms.read_at
FROM messages m
         LEFT JOIN messages_statuses ms ON
            m.message_id = ms.message_id AND
            ms.user_id = $2
WHERE m.message_id = $1;

-- name: ListRoomMessagesWithStatus :many
SELECT
    m.*,
    ms.is_read,
    ms.read_at
FROM messages m
         LEFT JOIN messages_statuses ms ON
            m.message_id = ms.message_id AND
            ms.user_id = $2
WHERE m.room_id = $1
ORDER BY m.created_at DESC
LIMIT $3 OFFSET $4;