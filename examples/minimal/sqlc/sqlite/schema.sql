CREATE TABLE author (
  author_id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL
);

CREATE INDEX authors_name_idx ON author(name);

CREATE TABLE book (
  book_id INTEGER PRIMARY KEY AUTOINCREMENT,
  author_id INTEGER NOT NULL,
  isbn TEXT NOT NULL DEFAULT '',
  book_type TEXT NOT NULL DEFAULT 'FICTION',
  title BLOB NOT NULL,
  -- UUID as binary
  yr INTEGER NOT NULL DEFAULT 2000,
  available DATETIME NOT NULL DEFAULT (datetime('now')),
  tags TEXT NOT NULL,
  FOREIGN KEY (author_id) REFERENCES author(author_id)
);