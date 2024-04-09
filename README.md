This is an xk6 extension for k6 that uses the amqp protocol.

Usage:

To build:

`make build`

To run:

First, you will have to make sure you've set the proper environment variables:

AMQP_TOPIC => the topic or endpoint you'd like to send your messages to.
e.g. "topic://test-topic"

AMQP_CONN_STRING => the address of the endpoint to send the messages to.
e.g. "amqp://admin:admin@localhost:5672/" 

For convenience I'd recommend modifying the `Makefile` in the root and substituting
your values in the command. Then you can simply use:

`make run`