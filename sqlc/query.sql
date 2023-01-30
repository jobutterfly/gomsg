-- name: GetBoardThreads :many
SELECT * FROM threads
WHERE board_id = ?
ORDER BY date ASC;

-- name: GetThreads :many
SELECT * FROM threads
ORDER BY date ASC
LIMIT ?;

-- name: GetThread :one
SELECT * FROM threads
WHERE thread_id = ?
LIMIT 1;

-- name: GetThreadReplies :many
SELECT * FROM replies
WHERE thread_id = ?
ORDER BY date ASC;

-- name: DeleteThread :execresult
DELETE FROM threads
WHERE thread_id = ?;

-- name: CreateThread :execresult
INSERT INTO threads(title, comment, date, board_id)
VALUES (?, ?, ?, ?);

-- name: CreateReply :execresult
INSERT INTO replies(comment, date, thread_id)
VALUES (?, ?, ?);

-- name: CountReplies :one
SELECT COUNT(*) FROM replies 
WHERE thread_id = ?;

-- name: CountThreads :one
SELECT COUNT(*) FROM threads;

-- name: GetOldestThread :one
SELECT * FROM threads 
WHERE board_id = ?
ORDER BY date ASC
LIMIT 1;






