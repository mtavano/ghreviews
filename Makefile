GOOS=linux

build-linux bl: clean
	@echo "[ghreview] Building for linux"
	@GOOS=linux \
	GOARCH=amd64 \
	go build -o main ./cmd/server/main.go

run r:
	@echo "[ghreview] Running ..."
	@go run ./cmd/server/main.go

migrate-up:
	@echo "[ghreview] Migrate up"
	@cd database/migrations go buildd -o goose *.go
	@goose $(DATABASE_DRIVER) $(DATABASE_URL) up

clean c:
	@echo "[CLEAN] Cleaning files..."
	@rm main || true

