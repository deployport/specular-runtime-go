package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestModulePathFromTrustedValues tests ModulePathFromTrustedValues
func TestModulePathFromTrustedValues(t *testing.T) {
	mp := ModulePathFromTrustedValues("ns", "name")
	require.Equal(t, "ns", mp.Namespace)
	require.Equal(t, "name", mp.Name)
}
