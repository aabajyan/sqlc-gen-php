CREATE TABLE author (
                         author_id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
                         name text NOT NULL
) ENGINE=InnoDB;

CREATE INDEX authors_name_idx ON author(name(255));

CREATE TABLE book (
                       book_id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
                       author_id integer NOT NULL,
                       isbn varchar(255) NOT NULL DEFAULT '',
                       book_type VARCHAR(100) NOT NULL DEFAULT 'FICTION',
                       title binary(16) NOT NULL COMMENT "UUID",
                       yr integer NOT NULL DEFAULT 2000,
                       available datetime NOT NULL DEFAULT NOW(),
                       tags text NOT NULL
    -- CONSTRAINT FOREIGN KEY (author_id) REFERENCES authors(author_id)
) ENGINE=InnoDB;

CREATE INDEX books_title_idx ON book(title(255), yr);

/*
CREATE FUNCTION say_hello(s text) RETURNS text
  DETERMINISTIC
  RETURN CONCAT('hello ', s);
*/
