integration-tests:
	docker-compose up -d 
	sleep 1
	go test -v ./test/integration/...
	docker-compose stop