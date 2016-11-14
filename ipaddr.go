package sockaddr

import (
	"fmt"
	"net"
	"strings"
)

// Constants for the sizes of IPv3, IPv4, and IPv6 address types.
const (
	IPv3len = 6
	IPv4len = 4
	IPv6len = 16
)

// IPAddr is a generic IP address interface for IPv4 and IPv6 addresses,
// networks, and socket endpoints.
type IPAddr interface {
	SockAddr
	AddressBinString() string
	AddressHexString() string
	Cmp(SockAddr) int
	CmpAddress(SockAddr) int
	CmpPort(SockAddr) int
	FirstUsable() IPAddr
	Host() IPAddr
	IPPort() IPPort
	LastUsable() IPAddr
	Maskbits() int
	NetIP() *net.IP
	NetIPMask() *net.IPMask
	NetIPNet() *net.IPNet
	Network() IPAddr
	Octets() []int
}

// IPPort is the type for an IP port number for the TCP and UDP IP transports.
type IPPort uint16

// IPPrefixLen is a typed integer representing the prefix length for a given
// IPAddr.
type IPPrefixLen byte

// ipAddrAttrMap is a map of the IPAddr type-specific attributes.
var ipAddrAttrMap map[AttrName]func(IPAddr) string
var ipAddrAttrs []AttrName

func init() {
	ipAddrInit()
}

// NewIPAddr creates a new IPAddr from a string.  Returns nil if the string is
// not an IPv4 or an IPv6 address.
func NewIPAddr(addr string) (IPAddr, error) {
	ipv4Addr, err := NewIPv4Addr(addr)
	if err == nil {
		return ipv4Addr, nil
	}

	ipv6Addr, err := NewIPv6Addr(addr)
	if err == nil {
		return ipv6Addr, nil
	}

	return nil, fmt.Errorf("invalid IPAddr %v", addr)
}

// IPAddrAttr returns a string representation of an attribute for the given
// IPAddr.
func IPAddrAttr(ip IPAddr, selector AttrName) string {
	fn, found := ipAddrAttrMap[selector]
	if !found {
		return ""
	}

	return fn(ip)
}

// IPAttrs returns a list of attributes supported by the IPAddr type
func IPAttrs() []AttrName {
	return ipAddrAttrs
}

// ipAddrInit is called once at init()
func ipAddrInit() {
	// Sorted for human readability
	ipAddrAttrs = []AttrName{
		"host",
		"address",
		"port",
		"netmask",
		"network",
		"mask_bits",
		"binary",
		"hex",
		"first_usable",
		"last_usable",
		"octets",
	}

	ipAddrAttrMap = map[AttrName]func(ip IPAddr) string{
		"address": func(ip IPAddr) string {
			return ip.NetIP().String()
		},
		"binary": func(ip IPAddr) string {
			return ip.AddressBinString()
		},
		"first_usable": func(ip IPAddr) string {
			return ip.FirstUsable().String()
		},
		"hex": func(ip IPAddr) string {
			return ip.AddressHexString()
		},
		"host": func(ip IPAddr) string {
			return ip.Host().String()
		},
		"last_usable": func(ip IPAddr) string {
			return ip.LastUsable().String()
		},
		"mask_bits": func(ip IPAddr) string {
			return fmt.Sprintf("%d", ip.Maskbits())
		},
		"netmask": func(ip IPAddr) string {
			return ip.NetIPMask().String()
		},
		"network": func(ip IPAddr) string {
			return ip.Network().String()
		},
		"octets": func(ip IPAddr) string {
			octets := ip.Octets()
			octetStrs := make([]string, 0, len(octets))
			for _, octet := range octets {
				octetStrs = append(octetStrs, fmt.Sprintf("%d", octet))
			}
			return strings.Join(octetStrs, " ")
		},
		"port": func(ip IPAddr) string {
			return fmt.Sprintf("%d", ip.IPPort())
		},
	}
}
