DOCKER_REGISTRY= localhost:5000
IMAGE= "$(DOCKER_REGISTRY)"/body-measurement
VERSION= latest
EXPOSED_PORT= "8080"

.PHONY: all

build:
	docker build --no-cache -t $(IMAGE):$(VERSION) .

push:
	docker push $(IMAGE):$(VERSION)

run:
	docker run --rm -p 8080:$(EXPOSED_PORT) $(IMAGE):$(VERSION) --log-level=debug

test:
	docker build --rm --target=tester .