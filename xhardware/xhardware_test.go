package xhardware

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetHardwareInfo(t *testing.T) {
	cpus, err := GetCpuInfo()
	fmt.Printf("cpus %v \n", cpus)
	require.NoError(t, err)
	require.NotNil(t, cpus)
	require.NotEmpty(t, cpus)

	nets, err := GetNetInfo()
	fmt.Printf("nets %v \n", nets)
	require.NoError(t, err)
	require.NotNil(t, nets)
	require.NotEmpty(t, nets)

	disk, err := GetDiskInfo()
	fmt.Printf("disk %v \n", disk)
	require.NoError(t, err)
	require.NotNil(t, disk)
	require.NotEmpty(t, disk)

	hwi, err := GetHardwareInfo()
	fmt.Printf("hardware info %v \n", hwi)
	require.NoError(t, err)
	require.NotNil(t, hwi)

	json, _ := json.Marshal(hwi)
	fmt.Printf("hardware info %v \n", string(json))
}
