BINARY_DIR := bin
CONFIG     := config.json
EVENTS     := events
IMAGE      := impulse

.PHONY: all test run build clean docker-build docker-run

all: test build

test:
	go test ./...

run:
	go run ./cmd/impulse $(CONFIG) $(EVENTS)

build:
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/impulse ./cmd/impulse

clean:
	rm -rf $(BINARY_DIR)

docker-build:
	docker build -t $(IMAGE) .

docker-run:
	docker run --rm $(IMAGE)
