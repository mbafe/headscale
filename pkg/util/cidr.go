package util

import (
	"fmt"
	"net/netip"
)

// SplitCIDR splits a prefix into two equal halves.
func SplitCIDR(prefix netip.Prefix) (netip.Prefix, netip.Prefix, error) {
	bits := prefix.Bits()
	total := prefix.Addr().BitLen()
	if bits >= total {
		return netip.Prefix{}, netip.Prefix{}, fmt.Errorf("cannot split a /%d prefix", bits)
	}

	first := netip.PrefixFrom(prefix.Addr(), bits+1)
	second, err := nextSubnet(prefix.Addr(), bits+1)
	if err != nil {
		return netip.Prefix{}, netip.Prefix{}, err
	}

	return first.Masked(), second.Masked(), nil
}

// nextSubnet returns the next subnet of the given prefix length after addr.
func nextSubnet(addr netip.Addr, bits int) (netip.Prefix, error) {
	p := netip.PrefixFrom(addr, bits).Masked()
	next := p.Addr()
	for i := 0; i < (addr.BitLen()-bits)+1; i++ {
		next = next.Next()
	}
	if !next.IsValid() {
		return netip.Prefix{}, fmt.Errorf("overflow computing next subnet")
	}
	return netip.PrefixFrom(next, bits), nil
}

// OverlapsCIDR reports whether two prefixes overlap.
func OverlapsCIDR(a, b netip.Prefix) bool {
	return a.Overlaps(b)
}

// SubtractCIDR returns the set of prefixes that cover 'whole' but not 'exclude'.
// This is a simplified implementation covering the common case.
// Note: when exclude equals or is broader than whole, an empty slice is returned
// rather than an error, which callers should treat as "nothing remains".
func SubtractCIDR(whole, exclude netip.Prefix) ([]netip.Prefix, error) {
	if !whole.Overlaps(exclude) {
		return []netip.Prefix{whole}, nil
	}
	if exclude.Bits() <= whole.Bits() {
		return []netip.Prefix{}, nil
	}

	var result []netip.Prefix
	current := whole
	for current.Bits() < exclude.Bits() {
		left, right, err := SplitCIDR(current)
		if err != nil {
			return nil, err
		}
		if exclude.Addr().Compare(right.Addr()) >= 0 && right.Contains(exclude.Addr()) {
			result = append(result, left)
			current = right
		} else {
			result = append(result, right)
			current = left
		}
	}
	return result, nil
}
