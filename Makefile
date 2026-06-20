.PHONY: build test clean run-learn run-generate run-enforce

BINARY := autoconfine
CMD := ./cmd/autoconfine

build:
	go build -o $(BINARY) $(CMD)

test:
	go test -race ./...

clean:
	rm -f $(BINARY) *.trace.jsonl *.seccomp.json coverage.out

run-learn: build
	./$(BINARY) learn --image nginx --duration 10s --out nginx.trace.jsonl

run-generate: build
	./$(BINARY) generate nginx.trace.jsonl --out nginx.seccomp.json

run-enforce: build
	./$(BINARY) enforce --profile nginx.seccomp.json -- podman run --rm nginx
