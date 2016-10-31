package factory

import (
	"fmt"

	"github.com/danielkrainas/cohesion/node"
)

var driverFactories = make(map[string]AgentFactory)

type AgentFactory interface {
	Create(parameters map[string]interface{}) (node.Agent, error)
}

func Register(name string, factory AgentFactory) {
	if factory == nil {
		panic("AgentFactory cannot be nil")
	}

	if _, registered := driverFactories[name]; registered {
		panic(fmt.Sprintf("AgentFactory named %s already registered", name))
	}

	driverFactories[name] = factory
}

func Create(name string, parameters map[string]interface{}) (node.Agent, error) {
	if factory, ok := driverFactories[name]; ok {
		return factory.Create(parameters)
	}

	return nil, InvalidAgentError{name}
}

type InvalidAgentError struct {
	Name string
}

func (err InvalidAgentError) Error() string {
	return fmt.Sprintf("Node agent not registered: %s", err.Name)
}
