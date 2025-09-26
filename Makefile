MAIN_PACKAGE := ./cmd/main.go
BINARY_NAME := aes-challenge-backend

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## build: Build the production code
.PHONY: build
build:
	@echo "Building binary for production..."
	@CGO_ENABLED=0 go build -o bin/${BINARY_NAME} ${MAIN_PACKAGE}

## build-dev: Build the development code
.PHONY: build-dev
build-dev:
	@echo "Building binary for development..."
	@CGO_ENABLED=0 go build -o bin/${BINARY_NAME} ${MAIN_PACKAGE}

## dev: Run the code development environment
.PHONY: dev
dev:
	@echo "Running in development environment..."
	@go run ${MAIN_PACKAGE}

## dev: Run the code prod environment
.PHONY: prod
prod:
	@echo "Running in production environment..."
	@go run ${MAIN_PACKAGE}

## watch: Run the application with reloading on file changes
.PHONY: watch
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

## update: Updates the packages and tidy the modfile
.PHONY: update
update:
	@go get -u ./...
	@go mod tidy -v

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: Format code and tidy modfile
.PHONY: tidy
tidy:
	@echo "Tidying up..."
	@go fmt ./...
	@go mod tidy -v

## test: Test the application
.PHONY: test
test:
	@echo "Testing..."
	@go test ./... -v

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: Print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
