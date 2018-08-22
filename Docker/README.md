# Wake-On-Lan Server

You can use Docker in both steps

## build
Docker will probably pull the `golang:latest` image (see [https://hub.docker.com/_/golang/])

    docker run --rm -u=$(id -u):$(id -g) -v "$PWD":/go -e CGO_ENABLED=0 golang go build src/wakeUp.go

You need to disable CGO to compile to a binary that can be used in the *scratch* image

## run
straight forward

	docker build -t wakeup .
	docker run -d -p 8000:8000 wakeup
