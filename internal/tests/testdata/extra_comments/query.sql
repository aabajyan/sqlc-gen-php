/* name: ListEntities :many */
/* @param bool|null $locked */
/* @param int|null $ownerId */
/* @param string|null $title */
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