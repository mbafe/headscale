package util

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIPPool(t *testing.T) {
	prefixes := []netip.Prefix{
		netip.MustParsePrefix("10.0.0.0/24"),
	}
	pool, err := NewIPPool(prefixes)
	require.NoError(t, err)
	assert.NotNil(t, pool)
}

func TestNewIPPoolInvalidPrefix(t *testing.T) {
	_, err := NewIPPool(nil)
	assert.Error(t, err)
}

func TestIPPoolNext(t *testing.T) {
	prefixes := []netip.Prefix{
		netip.MustParsePrefix("10.0.0.0/30"),
	}
	pool, err := NewIPPool(prefixes)
	require.NoError(t, err)

	// A /30 has 2 usable addresses (10.0.0.1 and 10.0.0.2);
	// network (10.0.0.0) and broadcast (10.0.0.3) are not allocated.
	ip1, err := pool.Next()
	require.NoError(t, err)
	assert.True(t, netip.MustParsePrefix("10.0.0.0/30").Contains(ip1))

	ip2, err := pool.Next()
	require.NoError(t, err)
	assert.True(t, netip.MustParsePrefix("10.0.0.0/30").Contains(ip2))
	assert.NotEqual(t, ip1, ip2)
}

func TestIPPoolExhausted(t *testing.T) {
	prefixes := []netip.Prefix{
		netip.MustParsePrefix("10.0.0.0/30"),
	}
	pool, err := NewIPPool(prefixes)
	require.NoError(t, err)

	// Exhaust the pool
	for {
		_, err := pool.Next()
		if err != nil {
			assert.ErrorContains(t, err, "exhausted")
			break
		}
	}
}

func TestIPPoolRelease(t *testing.T) {
	prefixes := []netip.Prefix{
		netip.MustParsePrefix("10.0.0.0/30"),
	}
	pool, err := NewIPPool(prefixes)
	require.NoError(t, err)

	ip, err := pool.Next()
	require.NoError(t, err)

	err = pool.Release(ip)
	require.NoError(t, err)

	// Released IPs should be re-issued before allocating new ones
	ip2, err := pool.Next()
	require.NoError(t, err)
	assert.Equal(t, ip, ip2)
}

func TestIPPoolReleaseUnknown(t *testing.T) {
	prefixes := []netip.Prefix{
		netip.MustParsePrefix("10.0.0.0/24"),
	}
	pool, err := NewIPPool(prefixes)
	require.NoError(t, err)

	// Releasing an IP outside the pool's prefixes should return an error
	err = pool.Release(netip.MustParseAddr("192.168.1.1"))
	assert.Error(t, err)
}

func TestIPPoolContains(t *testing.T) {
	prefixes := []netip.Prefix{
		netip.MustParsePrefix("10.0.0.0/24"),
	}
	pool, err := NewIPPool(prefixes)
	require.NoError(t, err)

	assert.True(t, pool.Contains(netip.MustParseAddr("10.0.0.1")))
	assert.False(t, pool.Contains(netip.MustParseAddr("192.168.1.1")))
}
