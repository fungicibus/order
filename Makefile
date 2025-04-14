.PHONY: generate-api-v1
generate-api-v1:
	oapi-codegen -package="v1" -generate types -o internal/api/v1/openapi_types.gen.go api/v1/openapi.yaml
	oapi-codegen -package="v1" -generate chi-server -o internal/api/v1/openapi_api.gen.go api/v1/openapi.yaml
	go mod tidy
