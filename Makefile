.PHONY: ama-go test

ama-go:
	go build ./cmd/ama-go

test: ama-go
	prove -e './ama-go' -v --ext '.ama' -r tests/tap
