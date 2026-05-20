test:
	go test -cover -coverprofile coverage.out -coverpkg "./internal/usecase/,./internal/delivery/handlers/" ./internal/tests/

cover:
	go tool cover -func coverage.out