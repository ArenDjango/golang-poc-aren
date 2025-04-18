Init:
go run cmd/main.go db init
go run cmd/main.go db migrate

Run http service:
go run cmd/main.go http

OpenAPI
generate -> swag init -g cmd/main.go --parseDependency --parseInternal
docs -> http://localhost:8087/swagger/index.html