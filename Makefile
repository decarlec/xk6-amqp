run:
	./k6 run test.js

build:
	xk6 build --with xk6-amqp=.

test:
	go run ./test/
