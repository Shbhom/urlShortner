create-migration:
	@if [ -z "$(name)" ]; then \
		echo "❌ Please provide a migration name using 'make create-migration name=your_migration_name'"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir internal/db/migrations $$name
