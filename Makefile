build:
	go build -o bin/notelify-users-service
	
serve-dev: build
	ENV=development ./bin/notelify-users-service

serve-dev-test: build
	ENV=development_test go test -v ./...

docker-push:
	docker build -t antonyinjila/notelify-users-service:latest --build-arg ENV=docker .
	docker push antonyinjila/notelify-users-service:latest

docker-run:
	docker run -p 8001:8001 ENV=docker antonyinjila/notelify-users-service:latest

docker-test:
	ENV=docker_test go test -v ./...





# dckr_pat_fVXnE212TTsa1qDzzyWXCmuSekQ