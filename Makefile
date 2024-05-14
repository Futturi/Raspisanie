sqlup:
	migrate -path migration -database "postgresql://root:root@localhost:5432/root?sslmode=disable" up
sqldown:
	migrate -path migration -database "postgresql://root:root@localhost:5432/root?sslmode=disable" down
.PHONY: sqlup, sqldown