
BINARY=sicra
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64 arm

deps:
	@echo "==> Updating build dependencies..."
	go get -u github.com/gocolly/colly/...

.PHONY: build
build:
	@echo "==> Building for all platforms..."
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -o build/$(BINARY)-$(GOOS)-$(GOARCH))))

.PHONY: deps build
