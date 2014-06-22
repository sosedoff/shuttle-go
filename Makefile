all: build

setup:
	go get

build:
	go build

clean:
	go clean
	rm ./shuttle-go