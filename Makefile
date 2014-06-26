all: build

deps:
	go get

build:
	go build

clean:
	go clean
	rm ./shuttle-go