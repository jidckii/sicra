
BINARY=sicra
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64

default: deps build
deps:
	@echo "==> Updating build dependencies..."
	go get -u github.com/gocolly/colly/...

build:
	@echo "==> Building for all platforms..."
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -o build/$(GOOS)/$(GOARCH)/$(BINARY))))

.PHONY: deps build
