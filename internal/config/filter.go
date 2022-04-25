package config

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)

// Filter represents a xdpfw filter policy
type Filter struct {
	MinLen           *int `json:"min_len"`   // Minimum frame length (ethernet header, IP header, L4 header, and data)
	MaxLen           *int `json:"max_len"`   // Maximum frame length (ethernet header, IP header, L4 header, and data)
	PacketsPerSecond *int `json:"pps"`       // Packets per second that a source IP can send before matching
	BytesPerSecond   *int `json:"bps"`       // Bytes per second that a source IP can send before matching
	BlockTime        *int `json:"blocktime"` // Time in seconds to block the source IP if the rule matches and the action is block (0). Default value is 1.

	// IP options
	TypeOfService *int    `json:"tos"`     // IP TOS field
	SrcIP         *string `json:"srcip"`   // Source IPv4 address
	DstIP         *string `json:"dstip"`   // Destination IPv4 address
	SrcIP6        *string `json:"srcip6"`  // Source IPv6 address
	DstIP6        *string `json:"dstip6"`  // Destination IPv6 address
	MinTTL        *int    `json:"min_ttl"` // Minimum TTL that the packet must match
	MaxTTL        *int    `json:"max_ttl"` // Maximum TTL that the packet must match

	// TCP Options
	TCPEnabled *bool `json:"tcp_enabled"` // Should TCP options be checked?
	TCPSrcPort *int  `json:"tcp_sport"`   // Source TCP port
	TCPDstPort *int  `json:"tcp_dport"`   // Destination TCP port
	TCPFlagURG *bool `json:"tcp_urg"`     // TCP URG flag
	TCPFlagACK *bool `json:"tcp_ack"`     // TCP ACK flag
	TCPFlagRST *bool `json:"tcp_rst"`     // TCP RST flag
	TCPFlagPSH *bool `json:"tcp_psh"`     // TCP PSH flag
	TCPFlagSYN *bool `json:"tcp_syn"`     // TCP SYN flag
	TCPFlagFIN *bool `json:"tcp_fin"`     // TCP FIN flag

	// UDP Options
	UDPEnabled *bool `json:"udp_enabled"` // Should UDP options be checked?
	UDPSrcPort *int  `json:"udp_sport"`   // Source UDP port
	UDPDstPort *int  `json:"udp_dport"`   // Destination UDP port

	// ICMP Options
	ICMPEnabled *bool `json:"icmp_enabled"` // Should ICMP options be checked?
	ICMPCode    *int  `json:"icmp_code"`    // ICMP code
	ICMPType    *int  `json:"icmp_type"`    // ICMP type
}

// String returns a string representation of the filter in xdpfw syntax
func (f *Filter) String() string {
	s := `{
  enabled = true,
  action = 0,
`

	v := reflect.ValueOf(f).Elem()
	vT := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsNil() {
			continue
		}
		value := v.Field(i).Elem().Interface()
		tag := vT.Field(i).Tag.Get("json")

		// Add quotes to represent empty strings
		if v.Field(i).Elem().Type() == reflect.TypeOf("") {
			value = "\"" + v.Field(i).Elem().String() + "\""
		}

		s += fmt.Sprintf("  %s = %v", tag, value)

		// Add a trialing comma if this is not the last field
		if v.NumField() != i+1 {
			s += ","
		}
		s += "\n"
	}

	return s + "}"
}

// Validate returns an error if the filter is invalid
func (f *Filter) Validate() error {
	if f.SrcIP != nil && net.ParseIP(*f.SrcIP) == nil {
		return fmt.Errorf("invalid source IP address: %s", *f.SrcIP)
	}

	if f.DstIP != nil && net.ParseIP(*f.DstIP) == nil {
		return fmt.Errorf("invalid destination IP address: %s", *f.DstIP)
	}

	if f.SrcIP6 != nil && net.ParseIP(*f.SrcIP6) == nil {
		return fmt.Errorf("invalid source IPv6 address: %s", *f.SrcIP6)
	}

	if f.DstIP6 != nil && net.ParseIP(*f.DstIP6) == nil {
		return fmt.Errorf("invalid destination IPv6 address: %s", *f.DstIP6)
	}

	return nil
}

// FromJSON converts a JSON string into a filter
func FromJSON(s string) (*Filter, error) {
	var f Filter
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return nil, err
	}

	return &f, nil
}
