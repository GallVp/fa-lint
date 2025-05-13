all: clean build bats

build:
	go build -ldflags "$$(<version)" -o bin/

clean:
	go mod tidy
	rm -rf bin

bats:
	bats -p test/test.bats

cgo:
	@$(MAKE) cgo-build OS=windows ARCH=amd64
	@$(MAKE) cgo-build OS=linux ARCH=amd64
	@$(MAKE) cgo-build OS=linux ARCH=arm64
	@$(MAKE) cgo-build OS=darwin ARCH=amd64
	@$(MAKE) cgo-build OS=darwin ARCH=arm64

# Helper target for CGO builds
cgo-build:
	@mkdir -p bin
	CGO_ENABLED=1 GOOS=$(OS) GOARCH=$(ARCH) \
	go build -ldflags "$$(<version)" -o bin/fa-lint-$(OS)-$(ARCH)/
	tar -czf bin/fa-lint-$(OS)-$(ARCH).tar.gz -C bin/fa-lint-$(OS)-$(ARCH) .
	rm -rf bin/fa-lint-$(OS)-$(ARCH)
