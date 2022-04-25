package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolicyFilterJSON(t *testing.T) {
	configString := `
{
  "pps": 10,
  "bps": 100,
  "blocktime": 1,
  "tos": 0,
  "srcip": "192.0.2.1"
}
`

	filter, err := FromJSON(configString)
	assert.Nil(t, err)
	assert.Equal(t, 10, *filter.PacketsPerSecond)
	assert.Equal(t, 100, *filter.BytesPerSecond)
	assert.Equal(t, 1, *filter.BlockTime)
	assert.Equal(t, 0, *filter.TypeOfService)
	assert.Equal(t, "192.0.2.1", *filter.SrcIP)
}
