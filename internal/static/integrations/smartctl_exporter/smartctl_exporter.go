// Package smartctl_exporter embeds https://github.com/prometheus-community/smartctl_exporter
package smartctl_exporter //nolint:golint

import (
	"errors"
	"strings"

	se "github.com/prometheus-community/smartctl_exporter"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/alloy/internal/static/integrations"
)

// DefaultConfig holds the default settings for the smartctl_exporter integration
var DefaultConfig = Config{
	SmartctlLoc:        "/usr/sbin/smartctl",
	SmartctlDevExclude: "",
	SmartctlInterval:   "60s",
}

// Config controls the smartctl_exporter integration.
type Config struct {
	SmartctlLoc        string `yaml:"smartctl.path,omitempty"`
	SmartctlDevExclude string `yaml:"smartctl.device-exclude,omitempty"`
	SmartctlInterval   string `yaml:"smartctl.interval,omitempty"`
}

// UnmarshalYAML implements yaml.Unmarshaler for Config
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig

	type plain Config
	return unmarshal((*plain)(c))
}

// Name returns the name of the integration this config is for.
func (c *Config) Name() string {
	return "smartctl_exporter"
}

// InstanceKey returns the smartctl location (SmartctlLoc).
func (c *Config) InstanceKey(agentKey string) (string, error) {
	return c.SmartctlLoc, nil
}

// NewIntegration converts the config into an integration instance.
func (c *Config) NewIntegration(logger log.Logger) (integrations.Integration, error) {
	return New(logger, c)
}

func init() {
	integrations.RegisterIntegration(&Config{})
}

// New creates a new smartctl_exporter integration. The integration scrapes metrics
// using the smartctl tool.
func New(logger log.Logger, c *Config) (integrations.Integration, error) {
	// Ensure that SmartctlInterval is in the correct format
	if !strings.HasSuffix(c.SmartctlInterval, "s") {
		err := errors.New("smartctl.interval must end with 's' to indicate seconds")
		level.Error(logger).Log("msg", "invalid scrape interval", "err", err)
		return nil, err
	}

	conf := &se.Config{
		ScrapeURI:    c.SmartctlLoc,
		HostOverride: c.SmartctlDevExclude,
	}

	// Create a new exporter using smartctl_exporter
	seExporter := se.NewExporter(logger, conf)

	return integrations.NewCollectorIntegration(
		c.Name(),
		integrations.WithCollectors(seExporter),
	), nil
}
