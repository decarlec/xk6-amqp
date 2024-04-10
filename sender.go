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

// Represents a two way client connection to an amqp server. The client can send and receive messages.
type sender struct {
	vu      modules.VU
	metrics amqpMetrics
	ctx     context.Context
	cancel  context.CancelFunc
	obj     *goja.Object
	conn    *amqp.Conn
	session *amqp.Session
	sender  *amqp.Sender
	uri     string
	topic   string
}

// Creates an amqp connection a message sender, and a message reciever.
func (a *AmqpAPI) sender(c goja.ConstructorCall) *goja.Object {
	uri := c.Argument(0).String()
	topic := c.Argument(1).String()

	rt := a.vu.Runtime()

	if uri == "" {
		common.Throw(rt, errors.New("uri is required"))
	}
	if topic == "" {
		common.Throw(rt, errors.New("topic is required"))
	}

	sender := sender{
		vu:      a.vu,
		obj:     rt.NewObject(),
		uri:     uri,
		topic:   topic,
		metrics: a.metrics,
	}
	sender.ctx, sender.cancel = context.WithCancel(a.vu.Context())

	must := func(err error) {
		if err != nil {
			common.Throw(rt, err)
		}
	}

	must(sender.obj.DefineDataProperty(
		"Connect", rt.ToValue(sender.Connect), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE))

	must(sender.obj.DefineDataProperty(
		"Disconnect", rt.ToValue(sender.Disconnect), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE))

	must(sender.obj.DefineDataProperty(
		"Send", rt.ToValue(sender.Send), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE))

	return sender.obj
}

func (s *sender) Connect() {
	rt := s.vu.Runtime()

	var err error
	s.conn, err = amqp.Dial(s.ctx, s.uri, nil)
	if err != nil {
		common.Throw(rt, err)
	}

	s.session, err = s.conn.NewSession(s.ctx, nil)
	if err != nil {
		common.Throw(rt, fmt.Errorf("failed to connect create amqp session: %v", err))
	}

	s.sender, err = s.session.NewSender(s.ctx, s.topic, nil)
	if err != nil {
		common.Throw(rt, fmt.Errorf("failed to create amqp sender: %v", err))
	}
}

func (sender *sender) Disconnect() {
	sender.conn.Close()
	sender.cancel()
}

func (s *sender) Send(message string) (err error) {
	rt := s.vu.Runtime()

	s.checkConnected()

	// report stats
	startedAt := time.Now()
	defer func() {
		now := time.Now()
		s.reportStats(s.metrics.SendMessageTiming, nil, now, metrics.D(now.Sub(startedAt)))
		s.reportStats(s.metrics.SentBytes, nil, now, float64(len(message)))
		if err != nil {
			s.reportStats(s.metrics.SendMessageErrors, nil, now, 1)
		} else {
			s.reportStats(s.metrics.SendMessage, nil, now, 1)
		}
	}()

	// send a message
	err = s.sender.Send(s.ctx, amqp.NewMessage([]byte(message)), nil)
	if err != nil {
		common.Throw(rt, fmt.Errorf("failed to send amqp message: %v", err))
	}

	return
}

func (s *sender) checkConnected() {
	if s == nil || s.conn == nil || s.session == nil || s.sender == nil {
		common.Throw(s.vu.Runtime(), ErrNotConnected)
	}
}

func (s *sender) reportStats(metric *metrics.Metric, tags map[string]string, now time.Time, value float64) {
	state := s.vu.State()
	ctx := s.vu.Context()
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
