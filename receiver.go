package amqp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/go-amqp"
	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
)

// Represents a two way receiver connection to an amqp server. The receiver can send and receive messages.
type receiver struct {
	vu       modules.VU
	ctx      context.Context
	cancel   context.CancelFunc
	obj      *goja.Object
	conn     *amqp.Conn
	session  *amqp.Session
	receiver *amqp.Receiver
	metrics  amqpMetrics
	uri      string
	topic    string
}

var ErrNotConnected = errors.New("not connected")

// Creates an amqp connection a message sender, and a message reciever.
func (a *AmqpAPI) receiver(c goja.ConstructorCall) *goja.Object {
	uri := c.Argument(0).String()
	topic := c.Argument(1).String()

	rt := a.vu.Runtime()

	if uri == "" {
		common.Throw(rt, errors.New("uri is required"))
	}
	if topic == "" {
		common.Throw(rt, errors.New("topic is required"))
	}

	receiver := receiver{
		vu:      a.vu,
		obj:     rt.NewObject(),
		uri:     uri,
		topic:   topic,
		metrics: a.metrics,
	}
	receiver.ctx, receiver.cancel = context.WithCancel(a.vu.Context())

	must := func(err error) {
		if err != nil {
			common.Throw(rt, err)
		}
	}

	must(receiver.obj.DefineDataProperty(
		"Connect", rt.ToValue(receiver.Connect), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE))

	must(receiver.obj.DefineDataProperty(
		"Disconnect", rt.ToValue(receiver.Disconnect), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE))

	must(receiver.obj.DefineDataProperty(
		"Receive", rt.ToValue(receiver.Receive), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE))

	return receiver.obj
}

func (receiver *receiver) Connect() {
	rt := receiver.vu.Runtime()

	var err error
	receiver.conn, err = amqp.Dial(receiver.ctx, receiver.uri, nil)
	if err != nil {
		common.Throw(rt, err)
	}

	receiver.session, err = receiver.conn.NewSession(receiver.ctx, nil)
	if err != nil {
		common.Throw(rt, fmt.Errorf("failed to connect create amqp session: %v", err))
	}

	receiver.receiver, err = receiver.session.NewReceiver(receiver.ctx, receiver.topic, nil)
	if err != nil {
		common.Throw(rt, fmt.Errorf("failed to create amqp receiver: %v", err))
	}
}

func (receiver *receiver) Disconnect() {
	receiver.conn.Close()
	receiver.cancel()
}

func (r *receiver) Receive() (err error) {
	rt := r.vu.Runtime()

	r.checkConnected()

	startedAt := time.Now()

	// receive the message
	msg, err := r.receiver.Receive(r.ctx, nil)
	if err != nil || msg == nil {
		common.Throw(rt, fmt.Errorf("failed to receive amqp message: %v", err))
	}

	err = r.receiver.AcceptMessage(r.ctx, msg)
	if err != nil {
		common.Throw(rt, fmt.Errorf("failed to accept amqp message: %v", err))
	}

	// report stats
	defer func() {
		now := time.Now()
		diff := now.Sub(startedAt)
		if diff.Nanoseconds() > 0 {
			r.reportStats(r.metrics.ReceiveMessageTiming, nil, now, metrics.D(diff))
		}

		r.reportStats(r.metrics.ReceivedBytes, nil, now, float64(len(msg.GetData())))
		if err != nil {
			r.reportStats(r.metrics.ReceiveMessageErrors, nil, now, 1)
		} else {
			r.reportStats(r.metrics.ReceiveMessage, nil, now, 1)
		}
	}()

	return
}

func (r *receiver) checkConnected() {
	if r == nil || r.conn == nil || r.session == nil || r.receiver == nil {
		common.Throw(r.vu.Runtime(), ErrNotConnected)
	}
}

func (r *receiver) reportStats(metric *metrics.Metric, tags map[string]string, now time.Time, value float64) {
	state := r.vu.State()
	ctx := r.vu.Context()
	// if state == nil || ctx == nil {
	// 	return
	// }

	metrics.PushIfNotDone(ctx, state.Samples, metrics.Sample{
		Time: now,
		TimeSeries: metrics.TimeSeries{
			Metric: metric,
			Tags:   metrics.NewRegistry().RootTagSet().WithTagsFromMap(tags),
		},
		Value: value,
	})
}
