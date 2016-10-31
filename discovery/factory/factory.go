package factory

import (
	"fmt"

	"github.com/danielkrainas/cohesion/discovery"
)

var strategyFactories = make(map[string]StrategyFactory)

type StrategyFactory interface {
	Create(parameters map[string]interface{}) (discovery.Strategy, error)
}

func Register(name string, factory StrategyFactory) {
	if factory == nil {
		panic("StrategyFactory cannot be nil")
	}

	if _, registered := strategyFactories[name]; registered {
		panic(fmt.Sprintf("StrategyFactory named %s already registered", name))
	}

	strategyFactories[name] = factory
}

func Create(name string, parameters map[string]interface{}) (discovery.Strategy, error) {
	if factory, ok := strategyFactories[name]; ok {
		return factory.Create(parameters)
	}

	return nil, InvalidStrategyError{name}
}

type InvalidStrategyError struct {
	Name string
}

func (err InvalidStrategyError) Error() string {
	return fmt.Sprintf("Strategy not registered: %s", err.Name)
}
