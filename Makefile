.PHONY: tail-watch tail-prod build push

build:
	docker build . --tag drywaters/bob

push:
	docker push drywaters/bob:latest

tail-watch: 
	tailwindcss -c ./tailwind/tailwind.config.js -i ./tailwind/styles.css -o ./static/styles.css --watch

tail-prod: 
	tailwindcss -c ./tailwind/tailwind.config.js -i ./tailwind/styles.css -o ./static/styles.css --minify
