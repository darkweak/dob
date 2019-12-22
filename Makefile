.PHONY: setup

build: ## Build project
	go build

setup: ## Setup and start project
	touch acme.json
	./dob
