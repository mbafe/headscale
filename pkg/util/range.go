package util

import (
	"fmt"
	"net/netip"
)

// IPRange represents an inclusive range of IP addresses.
type IPRange struct {
	First netip.Addr
	Last  netip.Addr
}

// NewIPRange creates a new IPRange from two addresses.
func NewIPRange(first, last netip.Addr) (IPRange, error) {
	if !first.IsValid() || !last.IsValid() {
		return IPRange{}, fmt.Errorf("invalid address")
	}
	if first.BitLen() != last.BitLen() {
		return IPRange{}, fmt.Errorf("address family mismatch")
	}
	if last.Compare(first) < 0 {
		return IPRange{}, fmt.Errorf("last address must be >= first address")
	}
	return IPRange{First: first, Last: last}, nil
}

// Contains reports whether addr is within the range.
func (r IPRange) Contains(addr netip.Addr) bool {
	return addr.Compare(r.First) >= 0 && addr.Compare(r.Last) <= 0
}

// ToPrefixes converts the IP range to a minimal set of CIDR prefixes.
func (r IPRange) ToPrefixes() []netip.Prefix {
	var prefixes []netip.Prefix
	current := r.First
	for current.Compare(r.Last) <= 0 {
		bits := current.BitLen()
		for bits > 0 {
			p := netip.PrefixFrom(current, bits-1).Masked()
			last, _ := LastUsableIP(p)
			if p.Addr().Compare(current) == 0 && last.Compare(r.Last) <= 0 {
				bits--
				continue
			}
			break
		}
		chosen := netip.PrefixFrom(current, bits)
		prefixes = append(prefixes, chosen)
		end, _ := LastUsableIP(chosen)
		if end.Compare(r.Last) >= 0 {
			break
		}
		current = end.Next()
		if !current.IsValid() {
			break
		}
	}
	return prefixes
}
