package version

import (
	"fmt"

	"github.com/danielkrainas/cohesion/cmd"
	"github.com/danielkrainas/cohesion/context"
)

func init() {
	cmd.Register("version", Info)
}

func run(ctx context.Context, args []string) error {
	fmt.Println("Cohesion v" + context.GetVersion(ctx))
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "version",
		Short: "`version`",
		Long:  "`version`",
		Run:   cmd.ExecutorFunc(run),
	}
)
