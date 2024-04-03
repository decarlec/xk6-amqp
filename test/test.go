package main

import (
	"context"

	"github.com/Azure/go-amqp"
)

// Testing out the amqp library
func main() {
	//Connect
	conn, err := amqp.Dial(context.TODO(), "amqp://user:password@domain", nil)
	if err != nil {
		panic(err)
	}

	//Create session
	session, err := conn.NewSession(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	//Create sender
	sender, err := session.NewSender(context.TODO(), "go-amqp", nil)
	if err != nil {
		panic(err)
	}

	// send a message
	err = sender.Send(context.TODO(), amqp.NewMessage([]byte("sent a message")), nil)
	if err != nil {
		panic(err)
	}

	// create a new receiver
	receiver, err := session.NewReceiver(context.TODO(), "go-amqp", nil)
	if err != nil {
		panic(err)
	}

	// receive the next message
	msg, err := receiver.Receive(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	err = receiver.AcceptMessage(context.TODO(), msg)
	if err == nil {
		println(string(msg.GetData()))
	}

	conn.Close()
}
