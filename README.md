# SQLC PHP DBAL Plugin

A SQLC plugin that generates PHP code using Doctrine DBAL for type-safe database operations. This plugin allows you to write SQL queries and automatically generates PHP classes with proper type hints and database abstraction.

## Features

- Generates PHP classes from SQL queries using Doctrine DBAL
- Supports MySQL (PostgreSQL support in development)
- Type-safe database operations
- Integration with Symfony framework
- Support for various SQL operations:
  - SELECT queries (single and multiple results)
  - INSERT operations
  - UPDATE operations
  - DELETE operations
  - Complex joins
  - Parameterized queries
  - Array parameters

## Installation

1. Download the latest release from the [releases page](https://github.com/lcarilla/sqlc-plugin-php-dbal/releases)
2. Add the plugin to your `sqlc.yaml` configuration:

```yaml
version: '2'
plugins:
- name: php
  wasm:
    url: https://github.com/lcarilla/sqlc-plugin-php-dbal/releases/download/v0.0.2/sqlc-gen-php.wasm
    sha256: 74f7a968592aeb6171113ad0cb972b7da9739c33f26738fbd6b2eee8893ce157
```

## Configuration

Configure your SQL queries in `sqlc.yaml`:

```yaml
sql:
- schema: sqlc/authors/mysql/schema.sql
  queries: sqlc/authors/mysql/query.sql
  engine: mysql
  codegen:
    - out: src/Sqlc/MySQL
      plugin: php
      options:
        package: "App\\Sqlc\\MySQL"
```

Exmaple of a complete `sqlc.yaml` config:
```yaml
version: '2'
plugins:
- name: php
  wasm:
    url: https://github.com/lcarilla/sqlc-plugin-php-dbal/releases/download/v0.0.2/sqlc-gen-php.wasm
    sha256: 74f7a968592aeb6171113ad0cb972b7da9739c33f26738fbd6b2eee8893ce157
sql:
- schema: sqlc/authors/mysql/schema.sql
  queries: sqlc/authors/mysql/query.sql
  engine: mysql
  codegen:
    - out: src/Sqlc/MySQL
      plugin: php
      options:
        package: "App\\Sqlc\\MySQL"
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

The plugin will generate PHP classes for the Models and Queries that you can use in your application:

```php
use App\Sqlc\MySQL\QueriesImpl;

$author = new QueriesImpl($connection)->getAuthor(authorId: 1);
```

## Development Status

- ✅ MySQL support
- 🚧 PostgreSQL support (Work in Progress)
- ✅ Basic CRUD operations
- ✅ Complex queries with joins
- ✅ Type-safe parameters
- ✅ Array parameters

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
