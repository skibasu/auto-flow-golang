coverage:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

coverage-html:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out