include app.testing.env
export $(shell sed 's/=.*//' app.testing.env)

# Function to load env variables from app.env
load-env:
	@export $(shell sed 's/=.*//' app.testing.env)

image:
	docker build -t simplebank:latest .

postgres: # run `make postgres` to create the postgres container on the postgres:latest image
	docker run --name postgres -p 5432:5432 -e GIN_MODE=release -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb: # run `make createdb` to create the simple_bank database
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb: # run `make dropdb` to drop the simple_bank database
	docker exec -it postgres dropdb simple_bank

migrateup: # run `make migrateup` to migrate up the database
	migrate -path db/migration -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose up

migrateup-aws-db:
	migrate -path db/migration -database "postgresql://$(POSTGRES_USER):$(POSTGRES_AWS_PASSWORD)@${POSTGRES_AWS_HOST}:$(POSTGRES_PORT)/$(POSTGRES_DB)" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose up 1

migratedown: # run `make migratedown` to migrate down the database
	migrate -path db/migration -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mocks/store.go github.com/santinofajardo/simpleBank/db/sqlc Store

proto:
 	rm -f pb/*.go \
	rm -f doc/swagger/*swagger.json \
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock proto evans redis
