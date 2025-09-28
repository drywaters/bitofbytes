.PHONY: docker-build docker-publish docker-push local run tail-watch tail-prod 

local:
	make -j 2 tail-watch run

build: tail-prod docker-build docker-push

docker-build:
	docker build . --tag 192.168.1.2:9000/bob:latest

docker-push:
	docker push 192.168.1.2:9000/bob:latest

docker-publish:
	make docker-build docker-push

run:
	air

tail-watch: 
	tailwindcss -i ./tailwind/styles.css -o ./static/styles.css --watch

tail-prod: 
	tailwindcss -i ./tailwind/styles.css -o ./static/styles.css --minify
