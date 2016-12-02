DIST_NAME := popuko
CONFIGURE_FILE := ./config.go

all: help

help:
	@echo "Specify the task"
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@exit 1

clean: ## Remove the exec binary.
	rm -rf ./$(DIST_NAME)

bootstrap:
	rm -rf vendor/
	go get -u github.com/mattn/gom
	gom install

new_config: $(CONFIGURE_FILE) ## Create the config file from our boilerplate.

build: $(DIST_NAME) ## Build the exec binary for youe machine.

build_linux_x64: ## Just an alias to build for some cloud instance.
	env GOOS=linux GOARCH=amd64 make build -C .

run: $(DIST_NAME) ## Execute the binary for youe machine.
	./$(DIST_NAME)

$(CONFIGURE_FILE):
	cp $@.example $@

$(DIST_NAME): clean $(CONFIGURE_FILE)
	go build -o $(DIST_NAME)

travis: bootstrap
	make build -C .
