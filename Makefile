run:
	./k6 run test.js

build:
	xk6 build --with xk6-amqp=github.com/decarlec/xk6-amqp@latest

test:
	go run ./test/
