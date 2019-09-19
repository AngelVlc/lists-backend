test:
	docker-compose run --rm app go test ./...

coverage:
	docker-compose run --rm app go test ./... -cover

coverage-html:
	docker-compose run --rm app go test ./... -cover -coverprofile coverage.out && go tool cover -html=coverage.out

fmt:
	go fmt . ./stores ./models

build:
	docker-compose build

up:
	docker-compose up
