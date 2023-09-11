build:
	go build -o bin/notelify-users-service

serve-prod: build
	./bin/notelify-users-service -env=prod

serve: build
	./bin/notelify-users-service -env=dev