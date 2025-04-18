Init:
go run cmd/main.go db init
go run cmd/main.go db migrate

Run http service:
go run cmd/main.go http

OpenAPI
generate -> swag init -g cmd/main.go --parseDependency --parseInternal
docs -> http://localhost:8087/swagger/index.html

As an example to consuming any third party IP I used https://ipinfo.io/
in user_service file I just created proxy for making request and show response to user