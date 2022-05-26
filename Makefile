integration-tests:
	docker-compose up -d 
	sleep 1
	aws --endpoint-url http://localhost:4566 s3 mb s3://organization
	go test -timeout 60s ./test/integration/...
	aws --endpoint-url http://localhost:4566 s3 rb --force s3://organization
	docker-compose stop