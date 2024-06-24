include app.env
export $(shell sed 's/=.*//' app.env)

# Function to load env variables from app.env
load-env:
	@export $(shell sed 's/=.*//' app.env)

postgres: # run `make postgres` to create the postgres container on the postgres:latest image
	docker run --name $(POSTGRES_CONTAINER_NAME) -p $(POSTGRES_PORT):$(POSTGRES_PORT) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -d $(POSTGRES_IMAGE)

createdb: # run `make createdb` to create the simple_bank database
	docker exec -it $(POSTGRES_CONTAINER_NAME) createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(POSTGRES_DB)

dropdb: # run `make dropdb` to drop the simple_bank database
	docker exec -it $(POSTGRES_CONTAINER_NAME) dropdb $(POSTGRES_DB)

migrateup: # run `make migrateup` to migrate up the database
	migrate -path db/migration -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose up

migratedown: # run `make migratedown` to migrate down the database
	migrate -path db/migration -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server
