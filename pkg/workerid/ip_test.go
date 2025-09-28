package workerid

import (
	"net"
	"testing"
)

func TestFromIP(t *testing.T) {
	ip := net.ParseIP("255.255.255.255")
	t.Log(FromIP(ip))
}
