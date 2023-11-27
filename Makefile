build:
	go build -o bin/notelify-users-service

serve-prod: build
	./bin/notelify-users-service -env=prod

serve-dev: build
	./bin/notelify-users-service -env=dev