# fhist Makefile

BINARY_NAME=fhist
INSTALL_DIR=$(HOME)/bin
ALT_INSTALL_DIR=$(HOME)/Utils

.PHONY: all build test clean install uninstall help

all: build test

## build: Compiles the fhist binary
build:
	go build -o $(BINARY_NAME) main.go

## test: Runs unit tests
test:
	go test -v ./...

## install: Compiles and installs fhist to ~/bin and ~/Utils
install: build
	mkdir -p $(INSTALL_DIR) $(ALT_INSTALL_DIR)
	killall $(BINARY_NAME) 2>/dev/null || true
	cp -f $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	cp -f $(BINARY_NAME) $(ALT_INSTALL_DIR)/$(BINARY_NAME)
	@echo "Successfully installed $(BINARY_NAME) to $(INSTALL_DIR) and $(ALT_INSTALL_DIR)"

## uninstall: Removes fhist binaries from install locations
uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	rm -f $(ALT_INSTALL_DIR)/$(BINARY_NAME)
	@echo "Removed $(BINARY_NAME) from $(INSTALL_DIR) and $(ALT_INSTALL_DIR)"

## clean: Removes built binary
clean:
	rm -f $(BINARY_NAME)
	go clean

## help: Shows available Makefile targets
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed -e 's/## //' | awk 'BEGIN {FS = ": "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
