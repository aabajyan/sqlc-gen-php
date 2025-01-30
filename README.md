inspired by https://github.com/sqlc-dev/sqlc-gen-kotlin
postgres support still WIP
## Usage:
```
version: '2'
plugins:
- name: php
  wasm:
    url: https://github.com/lcarilla/sqlc-plugin-php-dbal/releases/download/v0.0.1/sqlc-gen-kotlin.wasm
    checksum: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
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
