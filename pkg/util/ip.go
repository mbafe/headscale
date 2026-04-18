package util

import (
	"fmt"
	"net/netip"
)

// ParseCIDR parses a CIDR string and returns a netip.Prefix.
func ParseCIDR(cidr string) (netip.Prefix, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return netip.Prefix{}, fmt.Errorf("failed to parse CIDR %q: %w", cidr, err)
	}
	return prefix, nil
}

// ContainsIP reports whether the given prefix contains the given IP address.
func ContainsIP(prefix netip.Prefix, ip netip.Addr) bool {
	return prefix.Contains(ip)
}

// ParseAddr parses an IP address string and returns a netip.Addr.
func ParseAddr(addr string) (netip.Addr, error) {
	ip, err := netip.ParseAddr(addr)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("failed to parse IP address %q: %w", addr, err)
	}
	return ip, nil
}

// FirstUsableIP returns the first usable host IP in a prefix (network address + 1).
func FirstUsableIP(prefix netip.Prefix) (netip.Addr, error) {
	masked := prefix.Masked()
	ip := masked.Addr().Next()
	if !masked.Contains(ip) {
		return netip.Addr{}, fmt.Errorf("prefix %s has no usable addresses", prefix)
	}
	return ip, nil
}

// LastUsableIP returns the last usable host IP in a prefix.
// For IPv4, this is the address just before the broadcast address.
// Note: this iterates through the prefix, so avoid calling on very large IPv6 prefixes.
func LastUsableIP(prefix netip.Prefix) (netip.Addr, error) {
	masked := prefix.Masked()
	ip := masked.Addr()
	var last netip.Addr
	for next := ip.Next(); masked.Contains(next); next = next.Next() {
		last = next
	}
	if !last.IsValid() {
		return netip.Addr{}, fmt.Errorf("prefix %s has no usable addresses", prefix)
	}
	return last, nil
}
