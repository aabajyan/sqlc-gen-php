-- name: InsertAuthor :exec
-- @sqlc-param int $id
-- @sqlc-param string $name='hello'
-- @sqlc-param int $age
INSERT INTO
    author (id, name, age)
VALUES
    (:id, :name, :age);