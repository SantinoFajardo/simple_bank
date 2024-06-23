postgres: # run `make postgres` to create the postgres container on the postgres:latest image
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest 

createdb: # run `make createdb` to create the simple_bank database
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb: # run `make dropdb` to drop the simple_bank database
	docker exec -it postgres dropdb simple_bank

migrateup: # run `make migrateup` to migrate up the database
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown: # run `make migrateup` to migrate down the database
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test