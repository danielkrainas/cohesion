package agent

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/danielkrainas/cohesion/configuration"
	"github.com/danielkrainas/cohesion/context"
	"github.com/danielkrainas/cohesion/discovery"
	strategyFactory "github.com/danielkrainas/cohesion/discovery/factory"
	"github.com/danielkrainas/cohesion/node"
	nodeFactory "github.com/danielkrainas/cohesion/node/factory"
)

type Agent struct {
	context.Context

	config *configuration.Config

	discovery discovery.Strategy

	node node.Agent
}

func (agent *Agent) Run() error {
	context.GetLogger(agent).Info("agent started")
	defer context.GetLogger(agent).Info("agent stopped")

	agent.MonitorAndRecover()
	return nil
}

func (agent *Agent) MonitorAndRecover() {
	context.GetLogger(agent).Info("node monitor started")
	defer context.GetLogger(agent).Info("node monitor stopped")

	config := agent.config.Monitor
	sleepPeriod := time.Duration(config.Frequency)

	for {
		context.GetLogger(agent).Info("checking node status")
		if ok, err := agent.node.IsConnected(agent); !ok {
			context.GetLogger(agent).Info("node disconnected, starting discovery")

			addrs, err := agent.discovery.Locate(agent)
			if err != nil {
				context.GetLogger(agent).Errorf("error executing discovery: %v", err)
			} else if i, err = agent.node.Join(agent, addrs); err != nil {
				context.GetLogger(agent).Errorf("error joining discovery candidates: %v", err)
			} else {
				context.GetLogger(agent).Info("node joined: %d node(s)", i)
			}
		} else if err != nil {
			context.GetLogger(agent).Errorf("error checking node status: %v", err)
		} else {
			context.GetLogger(agent).Info("node ok")
		}

		context.GetLogger(agent).Info("node check completed")
		time.Sleep(sleepPeriod)
	}
}

func New(ctx context.Context, config *configuration.Config) (*Agent, error) {
	ctx, err := configureLogging(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error configuring logging: %v", err)
	}

	log := context.GetLogger(ctx)
	log.Info("initializing agent")

	d, err := configureDiscovery(config)
	if err != nil {
		return nil, err
	}

	n, err := configureNodeAgent(config)
	if err != nil {
		return nil, err
	}

	log.Infof("using %q logging formatter", config.Log.Formatter)
	log.Infof("using %q node agent", config.Node.Type())
	log.Infof("using %q discovery strategy", config.Discovery.Type())

	return &Agent{
		Context:   ctx,
		config:    config,
		discovery: d,
		node:      n,
	}, nil
}

func configureNodeAgent(config *configuration.Config) (node.Agent, error) {
	params := config.Node.Parameters()
	if params == nil {
		params = make(configuration.Parameters)
	}

	return nodeFactory.Create(config.Node.Type(), params)
}

func configureDiscovery(config *configuration.Config) (discovery.Strategy, error) {
	params := config.Discovery.Parameters()
	if params == nil {
		params = make(configuration.Parameters)
	}

	return strategyFactory.Create(config.Discovery.Type(), params)
}

func configureLogging(ctx context.Context, config *configuration.Config) (context.Context, error) {
	log.SetLevel(logLevel(config.Log.Level))
	formatter := config.Log.Formatter
	if formatter == "" {
		formatter = "text"
	}

	switch formatter {
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})

	case "text":
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
		})

	default:
		if config.Log.Formatter != "" {
			return ctx, fmt.Errorf("unsupported log formatter: %q", config.Log.Formatter)
		}
	}

	if len(config.Log.Fields) > 0 {
		var fields []interface{}
		for k := range config.Log.Fields {
			fields = append(fields, k)
		}

		ctx = context.WithValues(ctx, config.Log.Fields)
		ctx = context.WithLogger(ctx, context.GetLogger(ctx, fields...))
	}

	ctx = context.WithLogger(ctx, context.GetLogger(ctx))
	return ctx, nil
}

func logLevel(level configuration.LogLevel) log.Level {
	l, err := log.ParseLevel(string(level))
	if err != nil {
		l = log.InfoLevel
		log.Warnf("error parsing level %q: %v, using %q", level, err, l)
	}

	return l
}
