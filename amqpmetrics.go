package amqp

import (
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
)

type amqpMetrics struct {
	SentBytes     *metrics.Metric
	ReceivedBytes *metrics.Metric

	SendMessage       *metrics.Metric
	SendMessageTiming *metrics.Metric
	SendMessageErrors *metrics.Metric

	ReceiveMessage       *metrics.Metric
	ReceiveMessageTiming *metrics.Metric
	ReceiveMessageErrors *metrics.Metric
}

func registerMetrics(vu modules.VU) (amqpMetrics, error) {
	var err error
	var m amqpMetrics

	registry := vu.InitEnv().Registry

	m.SentBytes, err = registry.NewMetric(metrics.DataSentName, metrics.Counter, metrics.Data)
	if err != nil {
		return m, err
	}
	m.ReceivedBytes, err = registry.NewMetric(metrics.DataReceivedName, metrics.Counter, metrics.Data)
	if err != nil {
		return m, err
	}

	m.SendMessage, err = registry.NewMetric("amqp_messages_sent", metrics.Counter)
	if err != nil {
		return m, err
	}

	m.SendMessageTiming, err = registry.NewMetric("amqp_send_time", metrics.Trend, metrics.Time)
	if err != nil {
		return m, err
	}

	m.SendMessageErrors, err = registry.NewMetric("amqp_send_error_count", metrics.Counter)
	if err != nil {
		return m, err
	}

	m.ReceiveMessage, err = registry.NewMetric("amqp_messages_received", metrics.Counter)
	if err != nil {
		return m, err
	}

	m.ReceiveMessageTiming, err = registry.NewMetric("amqp_receive_time", metrics.Trend, metrics.Time)
	if err != nil {
		return m, err
	}

	m.ReceiveMessageErrors, err = registry.NewMetric("amqp_receive_error_count", metrics.Counter)
	if err != nil {
		return m, err
	}

	return m, nil
}
