.PHONY: rabbitmq_test test lint

rabbitmq_test: 
	go test ./test/rabbitmq_test.go -run ^TestRabbitMQExchanges -v -timeout 1m

test:
	go test ./... -v -timeout 1m

lint:
	golangci-lint run ./test
