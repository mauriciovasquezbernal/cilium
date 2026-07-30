// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package subnet

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	iputil "github.com/cilium/cilium/pkg/ip"
	cilium_api_v2alpha1 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2alpha1"
)

func prefix(t *testing.T, s string) iputil.Prefix {
	t.Helper()
	p, err := netip.ParsePrefix(s)
	require.NoError(t, err)
	return iputil.PrefixFrom(p)
}

// TestEntriesFromTopology verifies that a CiliumSubnetTopology decodes into
// positional subnet identities: one identity per subnet group, numbered from 1
// in spec order, shared by every CIDR of that group.
func TestEntriesFromTopology(t *testing.T) {
	cst := &cilium_api_v2alpha1.CiliumSubnetTopology{
		Spec: cilium_api_v2alpha1.CiliumSubnetTopologySpec{
			SubnetGroups: []cilium_api_v2alpha1.SubnetGroup{
				{
					Name:  "net-a",
					CIDRs: []iputil.Prefix{prefix(t, "10.0.0.1/24"), prefix(t, "10.10.0.1/24")},
				},
				{
					Name:  "net-b",
					CIDRs: []iputil.Prefix{prefix(t, "10.20.0.1/24")},
				},
				{
					Name:  "net-c",
					CIDRs: []iputil.Prefix{prefix(t, "2001:0db8:85a3::/64")},
				},
			},
		},
	}

	got, err := entriesFromTopology(cst)
	require.NoError(t, err)

	want := []topologyEntry{
		{prefix: netip.MustParsePrefix("10.0.0.1/24"), identity: 1},
		{prefix: netip.MustParsePrefix("10.10.0.1/24"), identity: 1},
		{prefix: netip.MustParsePrefix("10.20.0.1/24"), identity: 2},
		{prefix: netip.MustParsePrefix("2001:0db8:85a3::/64"), identity: 3},
	}
	assert.Equal(t, want, got)
}

// TestEntriesFromTopologyNil verifies a nil object yields no entries, which is
// what clears the table when the singleton is deleted.
func TestEntriesFromTopologyNil(t *testing.T) {
	got, err := entriesFromTopology(nil)
	require.NoError(t, err)
	assert.Empty(t, got)
}

// TestEntriesFromTopologyEmpty verifies an empty spec yields no entries.
func TestEntriesFromTopologyEmpty(t *testing.T) {
	got, err := entriesFromTopology(&cilium_api_v2alpha1.CiliumSubnetTopology{})
	require.NoError(t, err)
	assert.Empty(t, got)
}
