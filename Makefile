run:
	go run cmd/server.go

docker-build-image:
	docker build -t geovanedeveloper/rate-limit-api:latest -f Dockerfile .

docker-up:
	docker-compose up -d

docker-run:
	docker run --rm -p 8080:8080 geovanedeveloper/rate-limit-api:latest
	
test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

build-mocks:
	go install go.uber.org/mock/mockgen@latest
	~/go/bin/mockgen -source=pkg/ratelimit/ratelimit.go -destination=pkg/ratelimit/mock/ratelimit.go
	~/go/bin/mockgen -source=pkg/ratelimit/event.go -destination=pkg/ratelimit/mock/event.go