package main

import (
	"fmt"
	"reflect"
)

// Filter represents a xdpfw filter policy
type Filter struct {
	Enabled          bool `xdpfw:"enabled"`   // Should this rule be enabled?
	Allow            bool `xdpfw:"action"`    // Should the packet be allowed or blocked (0 = block, 1 = allow)
	MinLen           int  `xdpfw:"min_len"`   // Minimum frame length (ethernet header, IP header, L4 header, and data)
	MaxLen           int  `xdpfw:"max_len"`   // Maximum frame length (ethernet header, IP header, L4 header, and data)
	PacketsPerSecond int  `xdpfw:"pps"`       // Packets per second that a source IP can send before matching
	BytesPerSecond   int  `xdpfw:"bps"`       // Bytes per second that a source IP can send before matching
	BlockTime        int  `xdpfw:"blocktime"` // Time in seconds to block the source IP if the rule matches and the action is block (0). Default value is 1.

	// IP options
	TypeOfService int    `xdpfw:"tos"`     // IP TOS field
	SrcIP         string `xdpfw:"srcip"`   // Source IPv4 address
	DstIP         string `xdpfw:"dstip"`   // Destination IPv4 address
	SrcIP6        string `xdpfw:"srcip6"`  // Source IPv6 address
	DstIP6        string `xdpfw:"dstip6"`  // Destination IPv6 address
	MinTTL        int    `xdpfw:"min_ttl"` // Minimum TTL that the packet must match
	MaxTTL        int    `xdpfw:"max_ttl"` // Maximum TTL that the packet must match

	// TCP Options
	TCPEnabled bool `xdpfw:"tcp_enabled"` // Should TCP options be checked?
	TCPSrcPort int  `xdpfw:"tcp_sport"`   // Source TCP port
	TCPDstPort int  `xdpfw:"tcp_dport"`   // Destination TCP port
	TCPFlagURG bool `xdpfw:"tcp_urg"`     // TCP URG flag
	TCPFlagACK bool `xdpfw:"tcp_ack"`     // TCP ACK flag
	TCPFlagRST bool `xdpfw:"tcp_rst"`     // TCP RST flag
	TCPFlagPSH bool `xdpfw:"tcp_psh"`     // TCP PSH flag
	TCPFlagSYN bool `xdpfw:"tcp_psh"`     // TCP SYN flag
	TCPFlagFIN bool `xdpfw:"tcp_fin"`     // TCP FIN flag

	// UDP Options
	UDPEnabled bool `xdpfw:"udp_enabled"` // Should UDP options be checked?
	UDPSrcPort int  `xdpfw:"udp_sport"`   // Source UDP port
	UDPDstPort int  `xdpfw:"udp_dport"`   // Destination UDP port

	// ICMP Options
	ICMPEnabled bool `xdpfw:"icmp_enabled"` // Should ICMP options be checked?
	ICMPCode    int  `xdpfw:"icmp_code"`    // ICMP code
	ICMPType    int  `xdpfw:"icmp_type"`    // ICMP type
}

// String returns a string representation of the filter in xdpfw syntax
func (f *Filter) String() string {
	s := "{\n"

	v := reflect.ValueOf(f).Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Interface()
		tag := typeOfS.Field(i).Tag.Get("xdpfw")
		if v.Field(i).Type() == reflect.TypeOf("") && v.Field(i).String() == "" {
			value = "\"\""
		}

		s += fmt.Sprintf("  %s = %v,\n", tag, value)
	}

	return s + "}"
}

func main() {
	f := Filter{}
	fmt.Println(f.String())
}
