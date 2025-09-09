-- name: GetAuthor :one
SELECT
    *
FROM
    author
WHERE
    author_id = ?;

-- name: ListAuthors :many
SELECT
    *
FROM
    author
ORDER BY
    name;