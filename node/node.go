package node

import (
	"github.com/danielkrainas/cohesion/context"
)

type Agent interface {
	IsConnected(ctx context.Context) (bool, error)
	Join(ctx context.Context, addrs []string) (int, error)
}
