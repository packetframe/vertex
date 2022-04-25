# vertex

XDP based edge firewall

## Example Filters

##### Ratelimit SYN to 10 pps

```json
{"pps": 10, "tcp_syn": true}
```

##### Ratelimit IP to 1 kbps

```json
{"bps": 1024, "srcip": "192.0.2.1"}
``` 
