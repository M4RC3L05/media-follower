export CGO_ENABLED = 0

GO_DIRECT_DEPS := $(shell go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
GO_FLAGS = -trimpath -ldflags="-w -s"
CURRENT_GIT_TAG := $(shell git describe --tags --exact-match HEAD 2>/dev/null || echo "latest")

.DEFAULT_GOAL: help
.PHONY: help
help:
	@echo "Available targets:"
	@cat $(abspath $(lastword $(MAKEFILE_LIST))) | grep -oP '^[a-zA-Z_-]+(?=:)' | sort | xargs printf "  %s\n"

.PHONY:deps-update
deps-update:
	go get -u $(GO_DIRECT_DEPS)
	go mod tidy

.PHONY: code-check
code-check:
	go mod tidy --diff
	golangci-lint run ./...
	golangci-lint fmt --diff-colored ./...
	govulncheck -show verbose -test ./...

.PHONY: code-test
code-test:
	ginkgo -r -p --randomize-all --randomize-suites --fail-on-pending --fail-on-empty --keep-going --trace

.PHONY: code-compile-templates
code-compile-templates:
	templ generate

.PHONY: bundle-frontend-admin-server
bundle-frontend-admin-server:
	go run $(GO_FLAGS) cmd/build_frontend/main.go \
		./internal/entrypoints/admin_server/frontend/main.css \
		./internal/entrypoints/admin_server/frontend/favicon.ico \
		./internal/entrypoints/admin_server/frontend/.gitkeep \
		./internal/entrypoints/admin_server/dist

.PHONY: db-migrate
db-migrate:
	dbmate -u 'sqlite:./data/app.db' -d ./internal/storage/migrations -s ./internal/storage/schema.sql up
	go run $(GO_FLAGS) cmd/gen_db_types/main.go

.PHONY: entry-admin-server
entry-admin-server: code-compile-templates bundle-frontend-admin-server
	go run $(GO_FLAGS) cmd/main/main.go admin-server

.PHONY: entry-releases-feed-server
entry-releases-feed-server:
	go run $(GO_FLAGS) cmd/main/main.go releases-feed-server

.PHONY: entry-fetch-releases
entry-fetch-releases-ITUNES_MUSIC_RELEASES_PROVIDER:
	go run $(GO_FLAGS) cmd/main/main.go fetch-releases ITUNES_MUSIC_RELEASES_PROVIDER

.PHONY: docker-build-main
docker-build-main:
	docker build --platform linux/amd64,linux/arm64 -t docker.io/maingufu/media-follower:$(CURRENT_GIT_TAG) . -f ./.containers/containerfile --push

.PHONY: docker-build-migrator
docker-build-migrator:
	docker build --platform linux/amd64,linux/arm64 -t docker.io/maingufu/media-follower-migrator:$(CURRENT_GIT_TAG) . -f ./.containers/migrator/containerfile --push
