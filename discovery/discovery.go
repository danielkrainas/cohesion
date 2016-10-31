package discovery

import (
	"github.com/danielkrainas/cohesion/context"
)

type Strategy interface {
	Locate(ctx context.Context) ([]string, error)
}
