-- Create "author" table
CREATE TABLE `author` (
  `author_id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL
);
-- Create index "authors_name_idx" to table: "author"
CREATE INDEX `authors_name_idx` ON `author` (`name`);
-- Create "book" table
CREATE TABLE `book` (
  `book_id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `author_id` integer NOT NULL,
  `isbn` text NOT NULL DEFAULT '',
  `book_type` text NOT NULL DEFAULT 'FICTION',
  `title` blob NOT NULL,
  `yr` integer NOT NULL DEFAULT 2000,
  `available` datetime NOT NULL DEFAULT (datetime('now')),
  `tags` text NOT NULL,
  CONSTRAINT `0` FOREIGN KEY (`author_id`) REFERENCES `author` (`author_id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
