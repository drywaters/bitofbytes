.PHONY: docker-build local push run tail-watch tail-prod 

local:
	make -j 2 tail-watch run

docker-build:
	docker build . --tag drywaters/bob

run:
	go run ./cmd/bob/bob.go

push:
	docker push drywaters/bob:latest

tail-watch: 
	tailwindcss -c ./tailwind/tailwind.config.js -i ./tailwind/styles.css -o ./static/styles.css --watch

tail-prod: 
	tailwindcss -c ./tailwind/tailwind.config.js -i ./tailwind/styles.css -o ./static/styles.css --minify
