package xnet

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNet(t *testing.T) {
	lanAll := GetLanAll()

	require.NotEmpty(t, lanAll)

	ipv4, mac := GetLanFirst()
	require.NotEmpty(t, ipv4)
	require.NotEmpty(t, mac)
	t.Log(ipv4, mac)

	ipv4, mac = GetLanLast()
	require.NotEmpty(t, ipv4)
	require.NotEmpty(t, mac)
	t.Log(ipv4, mac)

	ipv4Regex := `^((25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)(\.|$)){4}$`
	matched := regexp.MustCompile(ipv4Regex).MatchString(ipv4)
	require.True(t, matched)

	ipv4Wan := GetWanIpv4()
	require.NotEmpty(t, ipv4Wan)
	t.Log(ipv4Wan)

	randomPort := GetAvailablePort(1)
	require.NotEqual(t, randomPort, 0)
	t.Log(randomPort)
}
