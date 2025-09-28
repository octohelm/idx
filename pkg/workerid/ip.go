package workerid

import (
	"net"
)

func FromIP(ip net.IP) uint32 {
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return uint32(ip[2])<<8 + uint32(ip[3])
}
