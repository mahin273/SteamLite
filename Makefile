.PHONY: all test run-auth clean

all: test

test:
	go test -v streamlite.local/...

run-auth:
	go run ./services/auth/cmd/server/main.go

clean:
	rm -rf bin/ tmp/ scratch/
