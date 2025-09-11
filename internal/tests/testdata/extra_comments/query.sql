/* name: ListEntities :many */
/* @sqlc-param bool|null $locked=null */
/* @sqlc-param int|null $owner_id=null */
/* @sqlc-param string|null $title=null */
SELECT
    *
FROM
    entity
WHERE
    (
        :locked IS NULL
        OR locked = :locked
    )
    AND (
        :owner_id IS NULL
        OR owner_id = :owner_id
    )
    AND (
        :title IS NULL
        OR title LIKE '%' || :title || '%'
    )
ORDER BY
    id DESC;