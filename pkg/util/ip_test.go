package util

import (
	"net/netip"
	"testing"
)

func TestParseCIDR(t *testing.T) {
	tests := []struct {
		cidr    string
		wantErr bool
	}{
		{"10.0.0.0/8", false},
		{"192.168.1.0/24", false},
		{"fd7a:115c:a1e0::/48", false},
		{"not-a-cidr", true},
		{"300.0.0.0/8", true},
		// edge case: host bits set should still parse as a valid prefix
		{"10.0.0.1/8", false},
	}
	for _, tt := range tests {
		_, err := ParseCIDR(tt.cidr)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseCIDR(%q) error = %v, wantErr %v", tt.cidr, err, tt.wantErr)
		}
	}
}

func TestContainsIP(t *testing.T) {
	prefix := netip.MustParsePrefix("10.0.0.0/8")
	if !ContainsIP(prefix, netip.MustParseAddr("10.1.2.3")) {
		t.Error("expected 10.1.2.3 to be in 10.0.0.0/8")
	}
	if ContainsIP(prefix, netip.MustParseAddr("192.168.1.1")) {
		t.Error("expected 192.168.1.1 to not be in 10.0.0.0/8")
	}
}

func TestParseAddr(t *testing.T) {
	_, err := ParseAddr("10.0.0.1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	_, err = ParseAddr("bad-ip")
	if err == nil {
		t.Error("expected error for bad IP")
	}
}

func TestFirstUsableIP(t *testing.T) {
	prefix := netip.MustParsePrefix("192.168.1.0/24")
	ip, err := FirstUsableIP(prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip.String() != "192.168.1.1" {
		t.Errorf("expected 192.168.1.1, got %s", ip)
	}
}

func TestLastUsableIP(t *testing.T) {
	prefix := netip.MustParsePrefix("192.168.1.0/30")
	ip, err := LastUsableIP(prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// /30 has 4 addresses: .0 (network), .1, .2, .3 (broadcast)
	// last usable is .2, not .3
	if ip.String() != "192.168.1.2" {
		t.Errorf("expected 192.168.1.2, got %s", ip)
	}
}
