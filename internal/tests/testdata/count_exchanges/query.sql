/* name: CountExchanges :one */
SELECT
    COUNT(*)
FROM
    exchange;

/* name: GetIdByName :one */
SELECT
    id
FROM
    exchange
WHERE
    name = ?1;

/* name: GetAllExchanges :many */
SELECT
    id,
    name,
    COUNT(*)
FROM
    exchange;