package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolicyConfigJSON(t *testing.T) {
	configString := `
{
  "interface": "eth0",
  "update_time": 15,
  "filters": [
	{
	  "pps": 10,
	  "bytes_per_second": 100,
	  "blocktime": 1,
	  "tos": 0,
	  "srcip": "192.0.2.1"
	},
	{
	  "pps": 10,
	  "bps": 100,
	  "blocktime": 1,
	  "tos": 0,
	  "srcip": "192.0.2.1"
	}
  ]
}
`

	var config Config
	err := json.Unmarshal([]byte(configString), &config)
	assert.Nil(t, err)

	expected := `filters = (
  {
    enabled = true,
    action = 0,
    pps = 10,
    blocktime = 1,
    tos = 0,
    srcip = "192.0.2.1",
  },
  {
    enabled = true,
    action = 0,
    pps = 10,
    bps = 100,
    blocktime = 1,
    tos = 0,
    srcip = "192.0.2.1",
  }
);
`

	assert.Equal(t, expected, config.String())
}
