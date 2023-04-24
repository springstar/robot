APP=robot
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
MAKEFILE_LIST=Makefile
PROTO_DIR=msg/protocol
STUBS_DIR=pb

.PHONY: build
build: clean

	go build -o ${robot} .

.PHONY: run
run:
	go run -race .

.PHONY: clean
clean:
	go clean		

.PHONY: help

.PHONY: proto
proto:
	rm -rf $(STUBS_DIR) 2>/dev/null
	mkdir $(STUBS_DIR)
	protoc --go_out=$(STUBS_DIR) \
	$(PROTO_DIR)/*.proto
		
	



help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'