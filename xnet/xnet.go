package xnet

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
)

func GetWanIpv4() string {
	providers := []string{
		"https://ifconfig.me/ip",
		"https://myip.ipip.net",
		"https://ipinfo.io",
	}
	for _, provider := range providers {
		resp, err := http.Get(provider)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])`)
		matches := re.FindAllString(string(body), -1)
		for _, ip := range matches {
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}
	return ""
}

func GetLanIpV4First() string {
	ip, _ := GetLanFirst()
	return ip
}

func GetLanIpV4Last() string {
	ip, _ := GetLanLast()
	return ip
}

// get ipv4 all of lan
func GetLanIpv4All() []string {
	returns := make([]string, 0, len(GetLanAll()))
	for k := range GetLanAll() {
		returns = append(returns, k)
	}
	return returns
}

func GetLanMacFirst() string {
	_, mac := GetLanFirst()
	return mac
}

func GetLanMacLast() string {
	_, mac := GetLanLast()
	return mac
}

// get ipmac of lan
func GetLanMacAll() []string {
	returns := make([]string, 0, len(GetLanAll()))
	for _, v := range GetLanAll() {
		returns = append(returns, v)
	}
	return returns
}

func GetLanFirst() (string, string) {
	all := GetLanAll()
	if len(all) == 0 {
		return "", ""
	}
	for ip, mac := range all {
		return ip, mac
	}
	return "", ""
}

func GetLanLast() (string, string) {
	all := GetLanAll()
	if len(all) == 0 {
		return "", ""
	}
	for ip, mac := range all {
		return ip, mac
	}
	return "", ""
}

// get ipv4 all of lan: map[ipv4]mac
func GetLanAll() map[string]string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	lanAll := make(map[string]string)
	for _, i := range interfaces {
		if !IsInterfacePhysical(i) {
			continue
		}
		if !IsInterfaceUp(i) {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			return nil
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() || ip.IsLinkLocalUnicast() ||
				ip.Equal(net.ParseIP("0.0.0.0")) || ip.Equal(net.ParseIP("127.0.0.1")) {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // 非ipv4
			}
			mac := i.HardwareAddr
			if mac == nil || mac.String() == "" {
				continue
			}
			lanAll[ip.String()] = mac.String()
		}
	}
	return lanAll
}

func IsInterfaceUp(inf net.Interface) bool {
	return inf.Flags&net.FlagUp != 0
}

func IsInterfacePhysical(inf net.Interface) bool {
	if inf.Flags&net.FlagLoopback != 0 {
		return false
	}
	if len(inf.HardwareAddr) == 0 {
		return false
	}
	if IsVirtualInterfaceName(inf.Name) {
		return false
	}
	return true
}

func IsVirtualInterfaceName(name string) bool {
	virtualKeywords := []string{
		// 通用虚拟网卡 / 容器环境 / 虚拟机环境
		"loopback", "virtual", "docker", "veth", "cali", "vbox", "vmnet", "vnic",
		"flannel", "macvlan", "weave", "virbr", "ovs", "wsl", "hyper-v", "eth0:avahi",

		// Linux 中常见的虚拟接口名
		"lo", "bond", "dummy", "ifb", "ipoib", "macvtap", "qbr", "qvb",
		"qvo", "qr-", "tap-", "tun", "tunl", "wg", "gretap", "gre", "ip6tnl",
		"br-", "br0", "br1", // Linux 桥接网络（br-xxx）和默认桥接名

		// macOS 虚拟接口名
		"awdl",   // Apple Wireless Direct Link
		"llw",    // Link-local only interface
		"ap",     // Apple Personal Area Network
		"bridge", // macOS 桥接接口（如 bridge0）
		"en5",    // 虚拟网卡常用编号之一（如 Parallels、VMware 创建）

		// Windows 中文系统中可能出现的名称
		"本地连接", "*",
	}
	for _, keyword := range virtualKeywords {
		if strings.Contains(strings.ToLower(name), keyword) {
			return true
		}
	}
	return false
}

func IsIpv4(ip string) bool {
	return net.ParseIP(ip).To4() != nil
}

func IsPrivate(ip string) bool {
	return net.ParseIP(ip).IsPrivate()
}

func IsGlobal(ip string) bool {
	return net.ParseIP(ip).IsGlobalUnicast()
}

func IsLinkLocal(ip string) bool {
	return net.ParseIP(ip).IsLinkLocalUnicast()
}

func IsLoopback(ip string) bool {
	return net.ParseIP(ip).IsLoopback()
}

func IsMulticast(ip string) bool {
	return net.ParseIP(ip).IsMulticast()
}

func IsInterfaceLocal(ip string) bool {
	return net.ParseIP(ip).IsInterfaceLocalMulticast()
}

func ListenTCP(addr string) (*net.TCPListener, error) {
	if addr == "" {
		addr = "0.0.0.0:0"
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func GetAvailablePort(retry int) int {
	port, _ := GetAvailablePortE(retry)
	return port
}

func GetAvailablePortE(retry int) (int, error) {
	if retry <= 0 {
		retry = 1
	}
	for i := 0; i < retry; i++ {
		l, err := ListenTCP("localhost:0")
		if err != nil {
			continue
		}
		port := l.Addr().(*net.TCPAddr).Port
		l.Close()
		return port, nil
	}
	return 0, nil
}

func IsPortAvailable(port int) bool {
	l, err := ListenTCP(fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return false
	}
	defer l.Close()
	return true
}
