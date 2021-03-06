package configuration

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"time"
)

type Frequency time.Duration

func (f *Frequency) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*f = Frequency(d)
	return nil
}

func (version *Version) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var versionString string
	err := unmarshal(&versionString)
	if err != nil {
		return err
	}

	newVersion := Version(versionString)
	if _, err := newVersion.major(); err != nil {
		return err
	}

	if _, err := newVersion.minor(); err != nil {
		return err
	}

	*version = newVersion
	return nil
}

type Parameters map[string]interface{}

type Driver map[string]Parameters

func (driver Driver) Type() string {
	var driverType []string

	for k := range driver {
		driverType = append(driverType, k)
	}

	if len(driverType) > 1 {
		panic("multiple drivers specified in the configuration or environment: %s" + strings.Join(driverType, ", "))
	}

	if len(driverType) == 1 {
		return driverType[0]
	}

	return ""
}

func (driver Driver) Parameters() Parameters {
	return driver[driver.Type()]
}

func (driver Driver) setParameter(key string, value interface{}) {
	driver[driver.Type()][key] = value
}

func (driver *Driver) UnmarshalText(text []byte) error {
	driverType := string(text)
	*driver = Driver{
		driverType: Parameters{},
	}

	return nil
}

func (driver *Driver) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var driverMap map[string]Parameters
	err := unmarshal(&driverMap)
	if err == nil && len(driverMap) > 0 {
		types := make([]string, 0, len(driverMap))
		for k := range driverMap {
			types = append(types, k)
		}

		if len(types) > 1 {
			return fmt.Errorf("Must provide exactly one driver type. provided: %v", types)
		}

		*driver = driverMap
		return nil
	}

	var driverType string
	if err = unmarshal(&driverType); err != nil {
		return err
	}

	*driver = Driver{
		driverType: Parameters{},
	}

	return nil
}

func (driver Driver) MarshalYAML() (interface{}, error) {
	if driver.Parameters() == nil {
		return driver.Type(), nil
	}

	return map[string]Parameters(driver), nil
}

type LogLevel string

func (logLevel *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var strLogLevel string
	err := unmarshal(&strLogLevel)
	if err != nil {
		return err
	}

	strLogLevel = strings.ToLower(strLogLevel)
	switch strLogLevel {
	case "error", "warn", "info", "debug":
	default:
		return fmt.Errorf("Invalid log level %s. Must be one of [error, warn, info, debug]", strLogLevel)
	}

	*logLevel = LogLevel(strLogLevel)
	return nil
}

type LogConfig struct {
	Level     LogLevel               `yaml:"level,omitempty"`
	Formatter string                 `yaml:"formatter,omitempty"`
	Fields    map[string]interface{} `yaml:"fields,omitempty"`
}

type MonitorConfig struct {
	Frequency Frequency `yaml:"frequency"`
}

type Config struct {
	Log       LogConfig     `yaml:"log"`
	Monitor   MonitorConfig `yaml:"monitor"`
	Node      Driver        `yaml:"node"`
	Discovery Driver        `yaml:"discovery"`
}

type v0_1Config Config

func newConfig() *Config {
	config := &Config{
		Log: LogConfig{
			Level:     "debug",
			Formatter: "text",
			Fields:    make(map[string]interface{}),
		},
	}

	return config
}

func Parse(rd io.Reader) (*Config, error) {
	in, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	p := NewParser("cohesion", []VersionedParseInfo{
		{
			Version: MajorMinorVersion(0, 1),
			ParseAs: reflect.TypeOf(v0_1Config{}),
			ConversionFunc: func(c interface{}) (interface{}, error) {
				if v0_1, ok := c.(*v0_1Config); ok {
					return (*Config)(v0_1), nil
				}

				return nil, fmt.Errorf("Expected *v0_1Config, received %#v", c)
			},
		},
	})

	config := new(Config)
	err = p.Parse(in, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
