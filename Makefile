.PHONY: check fmt vet test run

GATEWAY_DIR := gateway

# check is the gate to run before every commit: formatting, vet and tests with the race detector.
check:
	@test -z "$$(cd $(GATEWAY_DIR) && gofmt -l .)" || { echo "gofmt needed on:"; (cd $(GATEWAY_DIR) && gofmt -l .); exit 1; }
	cd $(GATEWAY_DIR) && go vet ./...
	cd $(GATEWAY_DIR) && go test -race ./...

fmt:
	cd $(GATEWAY_DIR) && gofmt -w .

vet:
	cd $(GATEWAY_DIR) && go vet ./...

test:
	cd $(GATEWAY_DIR) && go test -race ./...

run:
	cd $(GATEWAY_DIR) && go run ./cmd/gateway
