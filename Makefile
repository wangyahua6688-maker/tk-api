.PHONY: run run-bff run-business run-user tidy fmt test

run: run-bff

run-bff:
	go run . -f etc/tk-api.yaml

run-business:
	cd ../tk-business && go run . -f etc/business.yaml

run-user:
	cd ../tk-user && go run . -f etc/user.yaml

tidy:
	go mod tidy

fmt:
	gofmt -w ./

test:
	go test ./...
