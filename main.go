package main

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/danielkrainas/cohesion/cmd"
	_ "github.com/danielkrainas/cohesion/cmd/agent"
	"github.com/danielkrainas/cohesion/cmd/root"
	_ "github.com/danielkrainas/cohesion/cmd/version"
	"github.com/danielkrainas/cohesion/context"
	_ "github.com/danielkrainas/cohesion/discovery/echo"
	_ "github.com/danielkrainas/cohesion/node/serf"
)

var appVersion string

const DEFAULT_VERSION = "0.0.0-dev"

func main() {
	if appVersion == "" {
		appVersion = DEFAULT_VERSION
	}

	rand.Seed(time.Now().Unix())
	ctx := context.WithVersion(context.Background(), appVersion)

	dispatch := cmd.CreateDispatcher(ctx, root.Info)
	if err := dispatch(); err != nil {
		log.Fatalln(err)
	}
}
