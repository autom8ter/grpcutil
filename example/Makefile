.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "Makefile Commands:"
	@echo "----------------------------------------------------------------"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
	@echo "----------------------------------------------------------------"

.PHONY: proto
proto: ## regenerate gRPC code
	@echo "generating protobuf code..."
	@rm -rf gen
	@docker run -v `pwd`:/tmp colemanword/prototool:latest prototool generate
	@go fmt ./...

docker-build: ## build docker container containing build tools
	@docker build -t colemanword/prototool:latest .

docker-push: ## push docker container containing build tools
	@docker push colemanword/prototool:latest

run: ## run the server
	@go run main.go

descriptor:
	@docker run -v `pwd`:/tmp colemanword/prototool:latest prototool descriptor-set --include-imports --include-source-info >> descriptor.pb