package amqp

import (
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

// init is called by the Go runtime at application startup.
func init() {
	modules.Register("k6/x/amqp", New())
}

type (
	RootModule struct{}

	AmqpAPI struct {
		vu      modules.VU
		metrics amqpMetrics
	}
)

// NewModuleInstance implements the modules.Module interface and returns
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	m, err := registerMetrics(vu)
	if err != nil {
		common.Throw(vu.Runtime(), err)
	}
	return &AmqpAPI{vu: vu, metrics: m}
}

// Exports exposes the given object in ts
func (mi *AmqpAPI) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"Sender":   mi.sender,
			"Receiver": mi.receiver,
		},
	}
}

// Common options for creating senders and receivers
type Options struct {
	Uri   string
	Topic string
}

func New() *RootModule {
	return &RootModule{}
}

var (
	_ modules.Instance = &AmqpAPI{}
	_ modules.Module   = &RootModule{}
)
