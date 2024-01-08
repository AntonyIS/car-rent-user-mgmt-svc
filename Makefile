build:
	go build -o bin/notelify-users-service

serve: build
	./bin/notelify-users-service
	
 test:
	go test -v -tags=myenv ./...
	Env=dev go test -v -tags=myenv ./...


