package xpkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPkg(t *testing.T) {
	path, err := GetPkgPathE("github.com/spf13/cast")

	require.NoError(t, err)
	require.NotEmpty(t, path)

	path = GetPkgPath("github.com/spf13/cast")
	require.NotEmpty(t, path)
}
