.PHONY: run test fmt lint swagger docker clean

run:
	go run cmd/server/main.go

test:
	go test ./... -v -cover

fmt:
	gofmt -w .

lint:
	golangci-lint run ./...

swagger:
	swag init -g cmd/server/main.go -o docs

docker:
	docker compose up --build

clean:
	rm -rf docs/docs.go docs/swagger.json docs/swagger.yaml
	go clean