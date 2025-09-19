# SQLC PHP PDO Plugin (Fork)

> ⚠️ **Warning:** This project is in a very early stage. Expect bugs and breaking changes. Use with caution in production environments.

A fork of the original SQLC plugin, this version generates PHP code using native PDO for type-safe database operations. It allows you to write SQL queries and automatically generates PHP classes with proper type hints and database abstraction, without requiring Doctrine DBAL.

## Features

- Generates PHP classes from SQL queries using native PDO
- Supports MySQL and SQLite
- Type-safe database operations
- Support for various SQL operations:
  - SELECT queries (single and multiple results)
  - INSERT operations
  - UPDATE operations
  - DELETE operations
  - Complex joins
  - Parameterized queries
  - Array parameters

## Installation

1. Download the latest release from the [releases page](https://github.com/aabajyan/sqlc-gen-php/releases) of this fork
2. Add the plugin to your `sqlc.yaml` configuration:

```yaml
version: "2"
plugins:
  - name: php
    wasm:
      url: https://github.com/aabajyan/sqlc-gen-php/releases/download/v0.0.15/sqlc-gen-php.wasm
      sha256: 950c201f00b8dbe4045c09a206793cc413ff8133c1d5e512b9e57a20427b046a
```

## Configuration

Configure your SQL queries in `sqlc.yaml`:

```yaml
sql:
  - schema: sqlc/authors/mysql/schema.sql
    queries: sqlc/authors/mysql/query.sql
    engine: mysql # or sqlite
    codegen:
      - out: src/Sqlc/MySQL
        plugin: php
        options:
          package: "App\\Sqlc\\MySQL" # or your preferred namespace
```

Example of a complete `sqlc.yaml` config:

```yaml
version: '2'
plugins:
- name: php
  wasm:
    url: https://github.com/aabajyan/sqlc-gen-php/releases/download/v0.0.15/sqlc-gen-php.wasm
    sha256: 950c201f00b8dbe4045c09a206793cc413ff8133c1d5e512b9e57a20427b046a
sql:
- schema: sqlc/authors/mysql/schema.sql
  queries: sqlc/authors/mysql/query.sql
  engine: mysql # or sqlite
  codegen:
    - out: src/Sqlc/MySQL
  plugin: php
      options:
  package: "App\\Sqlc\\MySQL" # or your preferred namespace
```

### Options

- `package`: The PHP namespace for generated classes
- `out`: Output directory for generated code

## Example Usage

### Schema Definition

```sql
CREATE TABLE author (
    author_id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name text NOT NULL
) ENGINE=InnoDB;
```

### Query Definition

```sql
/* name: GetAuthor :one */
SELECT * FROM author
WHERE author_id = ?;
```

### Generated PHP Code

The plugin will generate PHP classes for the Models and Queries that you can use in your application. Usage with PDO:

```php
use App\Sqlc\MySQL\QueriesImpl;

$pdo = new PDO($dsn, $user, $password);
$author = (new QueriesImpl($pdo))->getAuthor(authorId: 1);
```

## Development Status

- ✅ MySQL support
- ✅ SQLite support
- 🚧 PostgreSQL support (Work in Progress)
- ✅ Basic CRUD operations
- ✅ Complex queries with joins
- ✅ Type-safe parameters
- ✅ Array parameters

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
