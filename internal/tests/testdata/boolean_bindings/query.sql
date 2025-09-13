-- name: setFlag :exec
INSERT INTO
    feature_flags (id, name, enabled)
VALUES
    (?, ?, ?);

-- name: listFlags :many
SELECT
    id,
    name,
    enabled
FROM
    feature_flags
ORDER BY
    id;