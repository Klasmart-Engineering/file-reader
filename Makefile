integration-tests:
	docker-compose up -d 
	sleep 1
	aws --endpoint-url http://localhost:4566 s3 ls 
	go test -timeout 60s ./test/integration/... -v
	docker-compose stop