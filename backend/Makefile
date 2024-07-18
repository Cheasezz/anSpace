test:
	go test -v -short -count=1 ./...

test100:
	go test -v -short -count=100 ./...

cover:
	go test -v -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out


test.integration: ### Run local integration tests
	docker run --name=test-db -e POSTGRES_PASSWORD=qwerty -p 5432:5432 -d --rm postgres

	GIN_MODE=release go test -v -count=1 ./integration-test/
	docker stop test-db

compose-up-integration-test: ### Run docker-compose with integration tests
	docker-compose -f docker-compose.integration.yml up --build --abort-on-container-exit --exit-code-from integration
	
gen:
	mockgen -source=internal/repository/psql/auth.go -destination=internal/repository/psql/mocks/mock_auth_repo.go 
	mockgen -source=internal/service/auth.go -destination=internal/service/mocks/mock_auth_service.go
	mockgen -source=pkg/auth/manager.go -destination=pkg/auth/mocks/mock_auth_manager.go
	mockgen -source=pkg/logger/logger.go -destination=pkg/logger/mocks/mock_logger.go
	mockgen -source=pkg/hasher/hasher.go -destination=pkg/hasher/mocks/mock_hasher.go