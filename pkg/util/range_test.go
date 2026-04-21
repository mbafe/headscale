package util

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIPRange(t *testing.T) {
	first := netip.MustParseAddr("10.0.0.1")
	last := netip.MustParseAddr("10.0.0.10")

	r, err := NewIPRange(first, last)
	require.NoError(t, err)
	assert.Equal(t, first, r.First)
	assert.Equal(t, last, r.Last)
}

func TestNewIPRangeInvalid(t *testing.T) {
	first := netip.MustParseAddr("10.0.0.10")
	last := netip.MustParseAddr("10.0.0.1")

	_, err := NewIPRange(first, last)
	require.Error(t, err)
}

// TestNewIPRangeEqual verifies that a range where first == last is valid (single-address range).
func TestNewIPRangeEqual(t *testing.T) {
	addr := netip.MustParseAddr("10.0.0.5")

	r, err := NewIPRange(addr, addr)
	require.NoError(t, err)
	assert.Equal(t, addr, r.First)
	assert.Equal(t, addr, r.Last)
	assert.True(t, r.Contains(addr))
}

func TestIPRangeContains(t *testing.T) {
	r, err := NewIPRange(
		netip.MustParseAddr("10.0.0.1"),
		netip.MustParseAddr("10.0.0.10"),
	)
	require.NoError(t, err)

	assert.True(t, r.Contains(netip.MustParseAddr("10.0.0.5")))
	assert.True(t, r.Contains(netip.MustParseAddr("10.0.0.1")))
	assert.True(t, r.Contains(netip.MustParseAddr("10.0.0.10")))
	assert.False(t, r.Contains(netip.MustParseAddr("10.0.0.11")))
	assert.False(t, r.Contains(netip.MustParseAddr("10.0.0.0")))
	// Also verify that an address from a completely different subnet is not contained.
	assert.False(t, r.Contains(netip.MustParseAddr("192.168.1.1")))
}

func TestIPRangeToPrefixes(t *testing.T) {
	r, err := NewIPRange(
		netip.MustParseAddr("10.0.0.0"),
		netip.MustParseAddr("10.0.0.255"),
	)
	require.NoError(t, err)

	prefixes := r.ToPrefixes()
	require.NotEmpty(t, prefixes)
	// Should produce a single /24
	assert.Equal(t, "10.0.0.0/24", prefixes[0].String())
}

// TestIPRangeToPrefixesUnaligned verifies that an unaligned range produces multiple prefixes.
func TestIPRangeToPrefixesUnaligned(t *testing.T) {
	r, err := NewIPRange(
		netip.MustParseAddr("10.0.0.1"),
		netip.MustParseAddr("10.0.0.10"),
	)
	require.NoError(t, err)

	prefixes := r.ToPrefixes()
	require.NotEmpty(t, prefixes)
	// An unaligned range should require more than one prefix to represent.
	assert.Greater(t, len(prefixes), 1)
}
