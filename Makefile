build:
	go build

clean:
	rm ./warlock

gotest:
	go test ./...

test: build gotest

