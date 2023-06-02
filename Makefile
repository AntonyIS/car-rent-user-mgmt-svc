build:
	go build -o bin/vp-user-mgmt-svc

serve: build
	./bin/vp-user-mgmt-svc