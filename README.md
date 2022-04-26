# vertex

XDP based edge firewall

## Example Filters

##### Ratelimit SYN to 10 pps

```json
{"tcp_enabled": true, "tcp_syn": true, "pps": 10}
```

##### Ratelimit IP to 1 kbps

```json
{"bps": 1024, "srcip": "192.0.2.1"}
``` 
