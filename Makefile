.PHONY: docker-build docker-publish docker-push local run tail-watch tail-prod

TAILWIND_VERSION ?= 4.1.13
TAILWIND_OS ?= linux
TAILWIND_ARCH ?= x64
TAILWIND_DOWNLOAD_URL := https://github.com/tailwindlabs/tailwindcss/releases/download/v$(TAILWIND_VERSION)/tailwindcss-$(TAILWIND_OS)-$(TAILWIND_ARCH)
TAILWIND_BIN := ./tailwind/bin/tailwindcss

local:
	make -j 2 tail-watch run

build: tail-prod docker-build docker-push

docker-build:
	docker build . --tag drywaters/bob

docker-push:
	docker push drywaters/bob:latest

docker-publish:
	make docker-build docker-push

run:
	air

tail-watch: tailwind-cli
	$(TAILWIND_BIN) -i ./tailwind/styles.css -o ./static/styles.css --watch

tail-prod: tailwind-cli
	$(TAILWIND_BIN) -i ./tailwind/styles.css -o ./static/styles.css --minify

tailwind-cli: $(TAILWIND_BIN)

$(TAILWIND_BIN):
	mkdir -p $(dir $(TAILWIND_BIN))
	curl -sSL $(TAILWIND_DOWNLOAD_URL) -o $(TAILWIND_BIN)
	chmod +x $(TAILWIND_BIN)
