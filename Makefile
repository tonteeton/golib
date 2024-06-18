SRC := $(wildcard *.go) $(wildcard */*.go)

all: format test audit

.PHONY: test
test:
	 ego-go test -v ./... -coverprofile=coverage.out


.PHONY: format
format:
	ego-go fmt ./...
	ego-go mod tidy -v

.PHONY: audit
audit:
	ego-go mod verify
	ego-go vet ./...

.PHONY: proto
proto: ereport/private_keys_report.pb.go

ereport/private_keys_report.pb.go: ereport/private_keys_report.proto
	protoc --go_out=paths=source_relative:./ -I. ereport/private_keys_report.proto


.PHONY: clean
clean:
	go clean
	rm -f coverage.out
