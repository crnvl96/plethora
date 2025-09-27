.PHONY: test cov

test:
	go test -cover ./...

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
