migrateup:
	./migrate -path ../db/migrations -database "mysql://root:p@ssw0rd@tcp(localhost:3306)/money_management?query" -verbose up
migratedown:
	./migrate -path ../db/migrations -database "mysql://root:p@ssw0rd@tcp(localhost:3306)/money_management?query" -verbose down

.PHONY: migrateup migratedown