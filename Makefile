default: deps build
deps:
	@echo "==> Updating build dependencies..."
	go get -u github.com/gocolly/colly/...

build:
	@echo "==> Building for all platforms..."
	GOOS=linux GOARCH=amd64 go build -o build/sicra
	tar -czf build/sicra-linux-amd64.tar.gz build/sicra

.PHONY: deps build
