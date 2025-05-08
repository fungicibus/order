.PHONY: generate-api-v1 lint goimports gofmt

generate-api-v1:
	oapi-codegen -package="v1" -generate types -o internal/api/v1/openapi_types.gen.go api/v1/openapi.yaml
	oapi-codegen -package="v1" -generate chi-server -o internal/api/v1/openapi_api.gen.go api/v1/openapi.yaml
	go mod tidy

# Run golangci-lint (https://github.com/golangci/golangci-lint).
lint:
	golangci-lint run ./internal/...

goimports:
	find . -type f -name "*.go" -exec goimports -local github.com/fungicibus -w {} \;

gofmt:
	find . -type f -name "*.go" -exec gofmt -w {} \;

up:
	docker compose --file docker/docker-compose.yml up
