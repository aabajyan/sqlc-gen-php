inspired by https://github.com/sqlc-dev/sqlc-gen-kotlin
postgres support still WIP
## Usage:
```
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
