all: clean build bats

build:
	go build -ldflags "$$(<version)" -o bin/

clean:
	go mod tidy

bats:
	bats -p test/test.bats