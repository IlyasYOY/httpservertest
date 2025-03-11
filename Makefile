.DEFAULT_GOAL = verify

.PHONY: lint 
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 \
		run --config=.golangci.yml \
		./...

.PHONY: test 
test:
	go test -shuffle=on -fullpath ./... -race 

.PHONY: build 
build:
	go build ./...

.PHONY: verify 
verify: lint test
