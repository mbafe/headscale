package util

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitCIDR(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		first  string
		second string
	}{
		{"ipv4 /24", "192.168.1.0/24", "192.168.1.0/25", "192.168.1.128/25"},
		{"ipv4 /16", "10.0.0.0/16", "10.0.0.0/17", "10.0.128.0/17"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := netip.MustParsePrefix(tt.input)
			first, second, err := SplitCIDR(p)
			require.NoError(t, err)
			assert.Equal(t, tt.first, first.String())
			assert.Equal(t, tt.second, second.String())
		})
	}
}

func TestSplitCIDRError(t *testing.T) {
	p := netip.MustParsePrefix("10.0.0.1/32")
	_, _, err := SplitCIDR(p)
	require.Error(t, err)
}

func TestOverlapsCIDR(t *testing.T) {
	a := netip.MustParsePrefix("10.0.0.0/8")
	b := netip.MustParsePrefix("10.1.0.0/16")
	c := netip.MustParsePrefix("192.168.0.0/16")

	assert.True(t, OverlapsCIDR(a, b))
	assert.False(t, OverlapsCIDR(a, c))
}

func TestSubtractCIDR(t *testing.T) {
	whole := netip.MustParsePrefix("10.0.0.0/24")
	exclude := netip.MustParsePrefix("10.0.0.0/25")

	result, err := SubtractCIDR(whole, exclude)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "10.0.0.128/25", result[0].String())
}

func TestSubtractCIDRNoOverlap(t *testing.T) {
	whole := netip.MustParsePrefix("10.0.0.0/24")
	exclude := netip.MustParsePrefix("192.168.0.0/24")

	result, err := SubtractCIDR(whole, exclude)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, whole, result[0])
}
