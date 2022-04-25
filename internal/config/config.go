package config

import (
	"errors"
	"strings"
)

type Config struct {
	Filters []*Filter `json:"filters"`
}

// String converts a Config into a xdpfw config file
func (c *Config) String() string {
	s := `filters = (
`

	for i := 0; i < len(c.Filters); i++ {
		filterStr := c.Filters[i].String()
		filterStr = strings.ReplaceAll(filterStr, "\n", "\n  ")
		filterStr = strings.ReplaceAll(filterStr, "{", "  {")

		s += filterStr
		if i != len(c.Filters)-1 {
			s += ","
		}
		s += "\n"
	}

	return s + ");\n"
}

// Validate checks that the config is valid
func (c *Config) Validate() error {
	if len(c.Filters) > 100 {
		return errors.New("max number of filters is 100 (eBPF jump limit)")
	}

	for i := 0; i < len(c.Filters); i++ {
		if err := c.Filters[i].Validate(); err != nil {
			return err
		}
	}

	return nil
}
