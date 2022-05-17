integration-tests:
	docker-compose up -d 
	sleep 1
	go test ./test/integration/...
	docker-compose stop