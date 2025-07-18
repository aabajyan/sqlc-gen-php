diff_migrations: ## Create versioned migration for production
	@atlas migrate diff schema \
		--dir 'file://migrations' \
		--to 'file://sqlc/sqlite/schema.sql' \
		--dev-url 'sqlite://file?mode=memory&_fk=1'

apply_migrations: ## Apply changes to database
	@atlas schema apply \
		--url 'sqlite://data.db' \
		--dev-url 'sqlite://dev?mode=memory' \
		--to 'file://sqlc/sqlite/schema.sql' 

rehash_migrations: ## Rehash
	@atlas migrate hash \
		--dir 'file://migrations'

make_migration: ## Make New Migration
	atlas migrate new --edit --dir 'file://migrations'