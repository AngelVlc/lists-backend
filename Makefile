test:
	go test ./...

run:
	go run .

coverage:
	go test ./... -cover

coverage-html:
	go test ./... -cover -coverprofile coverage.out && go tool cover -html=coverage.out

fmt:
	go fmt . ./stores ./models

docker-build:
	docker-compose build

docker-up:
	docker-compose up
