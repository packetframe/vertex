package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolicyFilterJSON(t *testing.T) {
	configString := `
{
  "interface": "eth0",
  "update_time": 15,
  "filters": [
	{
	  "action": 0,
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
	assert.Equal(t, "eth0", config.Interface)
	assert.Equal(t, 15, config.UpdateTime)
	assert.Equal(t, 1, len(config.Filters))
	assert.Equal(t, 0, *config.Filters[0].Action)
	assert.Equal(t, 10, *config.Filters[0].PacketsPerSecond)
	assert.Equal(t, 100, *config.Filters[0].BytesPerSecond)
	assert.Equal(t, 1, *config.Filters[0].BlockTime)
	assert.Equal(t, 0, *config.Filters[0].TypeOfService)
	assert.Equal(t, "192.0.2.1", *config.Filters[0].SrcIP)
}

func TestPolicyConfigJSON(t *testing.T) {
	configString := `
{
  "interface": "eth0",
  "update_time": 15,
  "filters": [
	{
	  "name": "test",
	  "action": 0,
	  "pps": 10,
	  "bytes_per_second": 100,
	  "blocktime": 1,
	  "tos": 0,
	  "srcip": "192.0.2.1"
	},
	{
	  "action": 1,
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

	expected := `interface = "eth0";
updatetime = 15;

filters = (
  {
    action = 0,
    pps = 10,
    blocktime = 1,
    tos = 0,
    srcip = "192.0.2.1",
  },
  {
    action = 1,
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
