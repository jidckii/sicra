PLATFORMS=darwin linux windows

default: deps build
deps:
	@echo "==> Updating build dependencies..."
	go get -u github.com/gocolly/colly/...

build:
	@echo "==> Building for all platforms..."
	$(foreach GOOS, $(PLATFORMS),\
	$(shell GOOS=$(GOOS) GOARCH=amd64 go build -o build/$(GOOS)/sicra && \
	tar -czf build/sicra-$(GOOS)-amd64.tar.gz build/$(GOOS)/sicra))

.PHONY: deps build
