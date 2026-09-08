// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package features

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeFailureExceptions(t *testing.T) {
	defaultExceptions := []string{"reason0", "reason1"}
	tests := []struct {
		inputExceptions    []string
		expectedExceptions []string
	}{
		// Empty list of reasons.
		{
			inputExceptions:    []string{},
			expectedExceptions: []string{},
		},
		// Add a reason to default list.
		{
			inputExceptions:    []string{"+reason2"},
			expectedExceptions: []string{"reason0", "reason1", "reason2"},
		},
		// Remove a reason from default list.
		{
			inputExceptions:    []string{"-reason1"},
			expectedExceptions: []string{"reason0"},
		},
		// Add a reason then remove it.
		{
			inputExceptions:    []string{"+reason2", "-reason2"},
			expectedExceptions: []string{"reason0", "reason1"},
		},
		// Remove a reason then add it back.
		{
			inputExceptions:    []string{"-reason1", "+reason1"},
			expectedExceptions: []string{"reason0", "reason1"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("InputExceptions: %v", test.inputExceptions), func(t *testing.T) {
			result := ComputeFailureExceptions(defaultExceptions, test.inputExceptions)

			// computeFailureExceptions doesn't guarantee the order of the
			// returned slice so we have to sort both slices.
			slices.Sort(result)
			assert.Equal(t, test.expectedExceptions, result)
		})
	}
}

func TestGetIPFamily(t *testing.T) {
	for _, tt := range []struct {
		name   string
		ip     string
		family IPFamily
	}{
		{
			name:   "empty",
			ip:     "",
			family: IPFamilyAny,
		},
		{
			name:   "invalid IPv4",
			ip:     "10.0.1.",
			family: IPFamilyAny,
		},
		{
			name:   "hostname",
			ip:     "example.com",
			family: IPFamilyAny,
		},
		{
			name:   "IPv4",
			ip:     "10.0.13.12",
			family: IPFamilyV4,
		},
		{
			name:   "IPv4 unspecified",
			ip:     "0.0.0.0",
			family: IPFamilyV4,
		},
		{
			name:   "IPv4 mapped IPv6",
			ip:     "::ffff:192.0.2.128",
			family: IPFamilyV4,
		},
		{
			name:   "IPv4 mapped IPv6 in hexadecimal",
			ip:     "::ffff:c000:280",
			family: IPFamilyV4,
		},
		{
			name:   "IPv4 compatible IPv6",
			ip:     "::192.0.2.128",
			family: IPFamilyV6,
		},
		{
			name:   "IPv6",
			ip:     "2001:db8::1",
			family: IPFamilyV6,
		},
		{
			name:   "IPv6 loopback",
			ip:     "::1",
			family: IPFamilyV6,
		},
		{
			name:   "IPv6 unspecified",
			ip:     "::",
			family: IPFamilyV6,
		},
		{
			name:   "IPv6 with zone",
			ip:     "fe80::1%eth0",
			family: IPFamilyV6,
		},
		{
			name:   "Prefix",
			ip:     "192.0.2.1/32",
			family: IPFamilyAny,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			family := GetIPFamily(tt.ip)
			if family != tt.family {
				t.Errorf("GetFamily(%q) = %v, want %v", tt.ip, family, tt.family)
			}
		})
	}
}
