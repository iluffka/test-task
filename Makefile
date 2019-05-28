.PHONY: run
run:
	$(info #App is running)
	go run cmd/main.go -env=$(env)

.PHONY: build-linux
build-linux:
	$(info #Building)
	cd cmd && env GOOS=linux go build -o ../bin/counter

.PHONY: test
test:
	$(info #Run tests...)
	go test -v ./internal/... ./cmd/...