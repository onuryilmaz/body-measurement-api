VERSION= latest
DATA_EXPOSED_PORT= "9092"
TRACKING_EXPOSED_PORT = "9093"

.PHONY: all

build:
	docker build --no-cache -t tracking-api:$(VERSION) --target=tracking-api .
	docker build --no-cache -t data-api:$(VERSION) --target=data-api .

run:
	docker run -d -p 9093:$(TRACKING_EXPOSED_PORT) --name tracking-api --net hpi tracking-api:$(VERSION) --log-level=debug
	docker run -d -p 9092:$(DATA_EXPOSED_PORT) --net hpi data-api:$(VERSION) --log-level=debug --tracking-address=http://tracking-api:9093/api/record


test:
	docker build --no-cache --rm --target=tester .

clean:
	rm cmd/data-api/data-api
	rm cmd/tracking-api/tracking-api
