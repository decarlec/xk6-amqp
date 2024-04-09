run:
	./k6 run test.js -e AMQP_TOPIC="topic://test-topic" -e AMQP_CONN_STRING="amqp://admin:admin@localhost:5672/"

build:
	xk6 build --with xk6-amqp=github.com/decarlec/xk6-amqp@latest

test:
	go run ./test/
