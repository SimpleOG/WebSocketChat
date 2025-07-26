-- name: CreateMessageStatus :one
INSERT INTO messages_statuses (message_id, user_id)
VALUES ($1, $2)
RETURNING *;


-- name: UpdateMessageStatus :one
UPDATE messages_statuses
SET
    is_read = $3,
    read_at = CASE WHEN $3 = TRUE THEN now() ELSE NULL END
WHERE message_id = $1 AND user_id = $2
RETURNING *;

-- name: GetMessageStatus :one
SELECT * FROM messages_statuses
WHERE message_id = $1 AND user_id = $2;

-- name: ListUnreadMessages :many
SELECT m.* FROM messages m
                    JOIN messages_statuses ms ON m.message_id = ms.message_id
WHERE ms.user_id = $1 AND ms.is_read = FALSE
ORDER BY m.created_at DESC;

-- name: CountUnreadMessages :one
SELECT COUNT(*) FROM messages_statuses
WHERE user_id = $1 AND is_read = FALSE;