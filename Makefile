build:
	go build -o bin/notlify-user-svc

serve: build
	./bin/notlify-user-svc