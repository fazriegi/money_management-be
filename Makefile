migrateup:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/dbname" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/dbname" -verbose down

.PHONY: migrateup migratedown