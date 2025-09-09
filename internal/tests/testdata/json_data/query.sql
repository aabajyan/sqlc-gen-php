-- name: GetAuthor :one
SELECT
    author_id,
    data
FROM
    author
WHERE
    author_id = ?;

-- name: ListAuthors :many
SELECT
    author_id,
    data
FROM
    author
ORDER BY
    json_extract(data, '$.name');

-- name: CreateAuthor :exec
INSERT INTO
    author (data)
VALUES
    (?);