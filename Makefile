all:
	go build -ldflags "$$(<version)" -o bin/

clean:
	go mod tidy