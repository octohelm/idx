package workerid

import (
	"net"
	"testing"
)

func TestFromIP(t *testing.T) {
	ip := net.ParseIP("255.255.255.255")
	if got := FromIP(ip); got != 65535 {
		t.Fatalf("expect 65535, got %d", got)
	}
}

func TestFromIP_Nil(t *testing.T) {
	if got := FromIP(nil); got != 0 {
		t.Fatalf("expect 0, got %d", got)
	}
}

func TestFromIP_IPv6(t *testing.T) {
	ip := net.ParseIP("2001:db8::1")
	if got := FromIP(ip); got != 0 {
		t.Fatalf("expect 0 for non-ipv4, got %d", got)
	}
}
