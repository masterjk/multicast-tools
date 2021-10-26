PROJECT_SRCDIR=github.com/masterjk/multicast-tools
VERSION="0.0.1-$(shell git rev-parse --short=8 HEAD)"
VERSION="latest"

DOCKER_IMAGE_RECEIVER=josephkiok/multicast-receiver:${VERSION}
DOCKER_IMAGE_SENDER=josephkiok/multicast-sender:${VERSION}

CONTAINER_GOLANG=golang:1.17.2

default: build

build:
	@echo ➭ Building images
	@docker build -t ${DOCKER_IMAGE_RECEIVER} -f Dockerfile.receiver .
	@docker build -t ${DOCKER_IMAGE_SENDER} -f Dockerfile.sender .

push:
	@echo ➭ Pushing Docker Images
	@docker push ${DOCKER_IMAGE_RECEIVER}
	@docker push ${DOCKER_IMAGE_SENDER}
