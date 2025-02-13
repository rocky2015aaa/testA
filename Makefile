# Define variables for reusability
IMAGE_NAME = shmoulana/hashocr
CONTAINER_NAME = hashocrsvr
PORT = 9081
VERSION := 1.0.0
BUILD := production
DATE := $(shell date +'%Y-%m-%d_%H:%M:%S')

# Check if Docker image exists
image_exists = $(shell docker images -q $(IMAGE_NAME):latest)

# Check if Docker container exists
container_exists = $(shell docker ps -aq -f name=$(CONTAINER_NAME))

# Target to build the Docker image if it doesn't already exist
build:
ifeq ($(image_exists),)
	@echo "Building Docker image: $(IMAGE_NAME):latest"
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		--no-cache \
		-t $(IMAGE_NAME):latest .
else
	@echo "Docker image $(IMAGE_NAME):latest already exists."
endif

# Target to run the Docker container if it's not already running
run:
ifeq ($(container_exists),)
	@echo "Running Docker container: $(CONTAINER_NAME)"
	docker run -d -p $(PORT):$(PORT) --name $(CONTAINER_NAME) --restart always $(IMAGE_NAME):latest
else
	@echo "Docker container $(CONTAINER_NAME) already exists."
endif

# Target to build and run only if it’s the first time (initial setup)
setup: build run

# Target to stop the container if it exists
stop:
ifeq ($(container_exists),)
	@echo "Docker container $(CONTAINER_NAME) does not exist."
else
	@echo "Stopping Docker container: $(CONTAINER_NAME)"
	docker stop $(CONTAINER_NAME)
endif

# Target to remove the container if it exists
clean:
ifeq ($(container_exists),)
	@echo "Docker container $(CONTAINER_NAME) does not exist."
else
	@echo "Removing Docker container: $(CONTAINER_NAME)"
	docker rm $(CONTAINER_NAME)
endif

# Target to remove the image if it exists
clean-image:
ifeq ($(image_exists),)
	@echo "Docker image $(IMAGE_NAME):latest does not exist."
else
	@echo "Removing Docker image: $(IMAGE_NAME):latest"
	docker rmi $(IMAGE_NAME):latest
endif

# Target to stop and remove container and image
clean-all: stop clean clean-image

# Target to rebuild and rerun everything with conditional checks
rebuild: clean-all
	@if [ -z "$(shell docker images -q $(IMAGE_NAME):latest)" ]; then \
		make build; \
	fi
	@if [ -z "$(shell docker ps -aq -f name=$(CONTAINER_NAME))" ]; then \
		make run; \
	fi

# .PHONY prevents targets from being mistaken for files
.PHONY: build run stop clean clean-image clean-all rebuild
