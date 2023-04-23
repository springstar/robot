APP=robot

.PHONY: build
build: clean

	go build -o ${robot} .

.PHONY: run
run:
	go run -race .

.PHONY: clean
clean:
	go clean		