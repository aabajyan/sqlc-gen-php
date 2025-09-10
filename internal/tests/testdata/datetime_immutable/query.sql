-- name: GetAuthorByCreatedAt :one
SELECT
    *
FROM
    author
WHERE
    created_at = ?;

-- name: AddAuthor :exec
INSERT INTO
    author (name, created_at)
VALUES
    (sqlc.arg(name), sqlc.narg(created_at));