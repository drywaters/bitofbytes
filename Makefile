.PHONY: tail-watch

tail-watch: 
	tailwindcss -c ./tailwind/tailwind.config.js -i ./tailwind/styles.css -o ./static/styles.css --watch
