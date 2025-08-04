run:
	@air

# Run Tailwind CSS in watch mode
tailwind:
	@npx tailwindcss -i views/css/styles.css -o public/styles.css --watch

# Run templ in watch mode with proxy to localhost:8080
templ:
	@templ generate -watch -proxy=http://localhost:8080

dev:
	@make -j2 tailwind templ run

build:
	go run worker_main.go
	mkdir -p dist
	cp -r public/* dist/
	cp -r static/* dist/
