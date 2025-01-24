/* name: GetAuthor :one */
SELECT * FROM author
WHERE id = ? LIMIT 1;

/* name: ListAuthors :many */
SELECT * FROM author
ORDER BY name;

/* name: CreateAuthor :execresult */
INSERT INTO author (
  name, bio
) VALUES (
  ?, ? 
);

/* name: DeleteAuthor :exec */
DELETE FROM author
WHERE id = ?;
