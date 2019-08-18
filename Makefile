test:
	go test ./... -cover -coverprofile=coverage.out

cover:
    go tool cover -func=coverage.out

cover-html:
    go tool cover -html=coverage.out