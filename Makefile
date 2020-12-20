POSTGRES_URL ?= postgres://apipath:apipath@localhost:54320/ghreviews?sslmode=disable
graphql_files = pkg/graph/exec.go pkg/graph/model.go

BIN := ghreviews-server
OUTBIN = bin/$(BIN)

default: run

$(graphql_files): gqlgen.yml
	go run github.com/99designs/gqlgen --verbose

graphql:
	go run github.com/99designs/gqlgen --verbose

build: graphql
	go build -o $(OUTBIN) github.com/mtavano/ghreviews/cmd/server

run: $(graphql_files)
	go run cmd/server/server.go

dev: $(graphql_files)
	go run cmd/server/server.go -db-verbose

clean:
	-rm -f $(graphql_files)
	-rm -f $(OUTBIN)

migrate:
	migrate -database $(POSTGRES_URL) -path migrations up

rollback:
	migrate -database $(POSTGRESQL_URL) -path migrations down 1

.PHONY: run build clean migrate graphql
