
help:                   ## This help dialog.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

install:                ## Locally fetch dependencies
	go get ./...

build:                  ## Build the local package
	go build cmd/main/main.go

build_release:          ## Build the local package in release mode
	go build cmd/main/main.go

run:                    ## Run the local package
	go run cmd/main/main.go

clean:                  ## Clean the build object and stop dependencies


fclean: clean           ## Reset the state of the package
	rm -f main

test:                   ## Execute the unit tests
	go test -v ./...

build_docker:           ## Build the docker image
	docker compose -f deployments/deployment.docker-compose.yml build

test_docker:            ## Run the test in a docker image
	docker compose -f test/unit-test.docker-compose.yml up --build --exit-code-from app-test
	docker compose -f test/unit-test.docker-compose.yml down -v --rmi local

run_docker:             ## Run the docker image and the dependency
	docker compose -f deployments/deployment.docker-compose.yml up

.PHONY: help install build build_release run clean fclean test build_docker test_docker run_docker
