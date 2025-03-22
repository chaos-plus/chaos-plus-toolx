package xhardware

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetHardwareInfo(t *testing.T) {
	cpus, err := GetCpuInfo()
	require.NoError(t, err)
	require.NotNil(t, cpus)
    require.NotEmpty(t, cpus)
    
	nets, err := GetHardwareInfo()
	require.NoError(t, err)
	require.NotNil(t, nets)
    require.NotEmpty(t, nets)
    
	disk, err := GetHardwareInfo()
	require.NoError(t, err)
	require.NotNil(t, disk)
    require.NotEmpty(t, disk)

	hwi, err := GetHardwareInfo()
	require.NoError(t, err)
	require.NotNil(t, hwi)

    json, _:= json.Marshal(hwi)
    fmt.Printf("hardware info %v \n", string(json))
}
