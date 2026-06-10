BINARY := peek

.PHONY: build install uninstall test clean dist

build:
	go build -o $(BINARY) .

# Cross-compile the same static binaries CI publishes to the rolling release.
dist:
	for target in linux/amd64 linux/arm64 darwin/arm64; do \
		GOOS=$${target%/*} GOARCH=$${target#*/} CGO_ENABLED=0 \
			go build -trimpath -ldflags='-s -w' -o "dist/$(BINARY)-$${target%/*}-$${target#*/}" . || exit 1; \
	done

# Installs into GOBIN (or GOPATH/bin, default ~/go/bin). No sudo needed.
install:
	go install .

uninstall:
	rm -f $$(go env GOPATH)/bin/$(BINARY)

test:
	go test ./...

clean:
	rm -rf $(BINARY) dist
