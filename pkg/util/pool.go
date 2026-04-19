package util

import (
	"fmt"
	"net/netip"
	"sync"
)

// IPPool manages allocation of IP addresses from a set of prefixes.
// It tracks which addresses have been allocated and provides
// thread-safe methods for acquiring and releasing addresses.
type IPPool struct {
	mu        sync.Mutex
	prefixes  []netip.Prefix
	allocated map[netip.Addr]struct{}
}

// NewIPPool creates a new IPPool from the given prefixes.
func NewIPPool(prefixes ...netip.Prefix) *IPPool {
	return &IPPool{
		prefixes:  prefixes,
		allocated: make(map[netip.Addr]struct{}),
	}
}

// Allocate returns the next available IP address from the pool.
// Returns an error if no addresses are available.
func (p *IPPool) Allocate() (netip.Addr, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, prefix := range p.prefixes {
		addr, err := FirstUsableIP(prefix)
		if err != nil {
			continue
		}
		last, err := LastUsableIP(prefix)
		if err != nil {
			continue
		}
		for ; addr.Compare(last) <= 0; addr = addr.Next() {
			if _, used := p.allocated[addr]; !used {
				p.allocated[addr] = struct{}{}
				return addr, nil
			}
		}
	}

	return netip.Addr{}, fmt.Errorf("no available addresses in pool")
}

// AllocateSpecific attempts to allocate a specific IP address.
// Returns an error if the address is already allocated or not in any pool prefix.
func (p *IPPool) AllocateSpecific(addr netip.Addr) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, used := p.allocated[addr]; used {
		return fmt.Errorf("address %s is already allocated", addr)
	}

	for _, prefix := range p.prefixes {
		if prefix.Contains(addr) {
			p.allocated[addr] = struct{}{}
			return nil
		}
	}

	return fmt.Errorf("address %s is not within any pool prefix", addr)
}

// Release marks an IP address as available for re-allocation.
// Returns an error if the address was not allocated.
func (p *IPPool) Release(addr netip.Addr) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, used := p.allocated[addr]; !used {
		return fmt.Errorf("address %s is not allocated", addr)
	}

	delete(p.allocated, addr)
	return nil
}

// IsAllocated reports whether the given address is currently allocated.
func (p *IPPool) IsAllocated(addr netip.Addr) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	_, used := p.allocated[addr]
	return used
}

// AllocatedCount returns the number of currently allocated addresses.
func (p *IPPool) AllocatedCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.allocated)
}

// Available returns the number of unallocated addresses across all prefixes.
// Note: for IPv4, network and broadcast addresses are excluded from the count.
// For host routes (/32 IPv4, /128 IPv6) the single address is considered usable.
func (p *IPPool) Available() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	count := 0
	for _, prefix := range p.prefixes {
		ones := prefix.Bits()
		total := prefix.Addr().BitLen() - ones
		if total <= 0 {
			// host route counts as 1 usable address
			count += 1
			continue
		}
		size := (1 << total)
		if !prefix.Addr().Is6() && total >= 2 {
			size -= 2
		}
		count += size
	}
	return count - len(p.allocated)
}
