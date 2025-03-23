package xcmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmd(t *testing.T) {
	version, err := RunWithResult("go", "version")
	require.NoError(t, err)
	require.Contains(t, version, "go version go")
}
