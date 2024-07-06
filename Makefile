.PHONY: tail-watch tail-prod

tail-watch: 
	tailwindcss -c ./tailwind/tailwind.config.js -i ./tailwind/styles.css -o ./static/styles.css --watch

tail-prod: 
	tailwindcss -c ./tailwind/tailwind.config.js -i ./tailwind/styles.css -o ./static/styles.css --minify
