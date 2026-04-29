// Package workerid 提供 worker ID 推导相关的辅助能力。
package workerid

import (
	"net"
)

// FromIP 从 IPv4 地址的最后两个字节推导出一个 worker ID。
// 当传入为空或不是 IPv4 地址时返回 0。
func FromIP(ip net.IP) uint32 {
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return uint32(ip[2])<<8 + uint32(ip[3])
}
