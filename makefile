.PYONY: arm linux build

svcpod-build: 
	@go build cmd/ -o svcPodKmsExporter

darwin:
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o svcPodKmsExporter

linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o svcPodKmsExporter

run:
	@go run ./

help:
	@echo "usage: make <option>"
	@echo "    help   : Show help"
	@echo "    build  : Build the binary of this project for current platform"
	@echo "    run	  : run the  project"
	@echo "    linux  : Build the amd64 linux binary of this project"
	@echo "    darwin : Build the arm64 darwin binary of this project"