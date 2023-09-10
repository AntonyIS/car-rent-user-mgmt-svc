build:
	go build -o bin/notelify-user-service

serve-prod: build
	./bin/notelify-user-service -env=prod

serve: build
	./bin/notelify-user-service -env=dev