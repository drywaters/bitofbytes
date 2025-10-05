.PHONY: configure-image build-github docker-build-push-github docker-build docker-publish docker-push ensure-image-tag local run tail-watch tail-prod

configure-image:
	$(eval REGISTRY ?= registry.bitofbytes.io)
	$(eval IMAGE_NAME ?= $(REGISTRY)/bob)
	$(eval SHORT_SHA := $(shell git rev-parse --short HEAD))
	$(eval IMAGE_TAG ?= $(SHORT_SHA))
	$(eval IMAGE := $(IMAGE_NAME):$(IMAGE_TAG))
	@true

ensure-image-tag: configure-image
	@test -n "$(strip $(SHORT_SHA))" || (echo "Unable to determine git short SHA. Ensure this is a git repository with at least one commit." >&2; exit 1)

run:
	air

local:
	make -j 2 tail-watch run

build: tail-prod docker-build docker-push

docker-build: ensure-image-tag
	docker build -f Docker/Dockerfile . --tag $(IMAGE)

docker-push: ensure-image-tag
	docker push $(IMAGE)

docker-publish:
	make docker-build docker-push

tail-watch:
	tailwindcss -i ./tailwind/styles.css -o ./static/styles.css --watch

tail-prod:
	tailwindcss -i ./tailwind/styles.css -o ./static/styles.css --minify

build-github: tail-prod configure-image docker-build-push-github

docker-build-push-github:
	echo ">> Building and pushing $(IMAGE)"
	-docker buildx inspect >/dev/null 2>&1 || docker buildx create --use
	docker buildx build -f Docker/Dockerfile . \
		--platform=linux/amd64 \
		-t $(IMAGE) \
		--push
