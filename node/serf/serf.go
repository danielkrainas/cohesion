package serf

import (
	"errors"

	serfClient "github.com/hashicorp/serf/client"

	"github.com/danielkrainas/cohesion/context"
	"github.com/danielkrainas/cohesion/node"
	"github.com/danielkrainas/cohesion/node/factory"
)

type driverFactory struct{}

func (f *driverFactory) Create(parameters map[string]interface{}) (node.Agent, error) {
	addr, ok := parameters["addr"].(string)
	if !ok || addr == "" {
		return nil, errors.New("configuration for `addr` missing or invalid")
	}

	c, err := serfClient.NewRPCClient(addr)
	if err != nil {
		return nil, err
	}

	return &driver{
		client: c,
	}, nil
}

func init() {
	factory.Register("serf", &driverFactory{})
}

type driver struct {
	client *serfClient.RPCClient
}

func (d *driver) IsConnected(ctx context.Context) (bool, error) {
	m, err := d.client.Members()
	return len(m) > 1, err
}

func (d *driver) Join(ctx context.Context, addrs []string) (int, error) {
	return d.client.Join(addrs, false)
}
