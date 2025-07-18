/* name: GetAuthor :one */
SELECT
    *
FROM
    author
WHERE
    author_id = ?;

/* name: ListAuthors :many */
SELECT
    *
FROM
    author
ORDER BY
    name;

/* name: GetBook :one */
SELECT
    *
FROM
    book
WHERE
    book_id = ?;

/* name: DeleteBook :exec */
DELETE FROM
    book
WHERE
    book_id = ?;

/* name: bookByTitleYear :many */
SELECT
    *
FROM
    book
WHERE
    title = UUID_TO_BIN(?)
    AND yr = ?;

/* name: bookByTags :many */
SELECT
    book_id,
    title,
    name,
    isbn,
    tags
FROM
    book
    LEFT JOIN author ON book.author_id = author.author_id
WHERE
    tags = ?;

/* name: bookByTagsMultiple :many */
SELECT
    book_id,
    title,
    name,
    isbn,
    tags
FROM
    book
    LEFT JOIN author ON book.author_id = author.author_id
WHERE
    tags IN (sqlc.slice(tags));

/* name: CreateAuthor :execresult */
INSERT INTO
    author (name)
VALUES
    (?);

/* name: CreateBook :execresult */
INSERT INTO
    book (
        author_id,
        isbn,
        book_type,
        title,
        yr,
        available,
        tags
    )
VALUES
    (
        ?,
        ?,
        ?,
        UUID_TO_BIN(?),
        ?,
        ?,
        ?
    );

/* name: UpdateBook :exec */
UPDATE
    book
SET
    title = ?,
    tags = ?
WHERE
    book_id = ?;

/* name: UpdateBookISBN :exec */
UPDATE
    book
SET
    title = ?,
    tags = ?,
    isbn = ?
WHERE
    book_id = ?;

/* name: DeleteAuthorBeforeYear :exec */
DELETE FROM
    book
WHERE
    yr < ?
    AND author_id = ?;

-- WHERE yr < sqlc.arg(min_publish_year) AND author_id = ?;