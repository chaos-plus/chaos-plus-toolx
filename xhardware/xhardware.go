package xhardware

import (
	"github.com/chaos-plus/chaos-plus-toolx/xnet"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/net"
)

type CpuInfo struct {
	Index      int32   `json:"index"`
	CoreNum    int32   `json:"coreNum"`
	VemdorId   string  `json:"vendorId"`
	PhysicalId string  `json:"physicalId"`
	ModelName  string  `json:"modelName"`
	MHZ        float64 `json:"mhz"`
}

type NetInfo struct {
	Index int32  `json:"index"`
	Name  string `json:"name"`
	MTU   int32  `json:"mtu"`
	MAC   string `json:"mac"`
}

type DiskInfo struct {
	Index  int32  `json:"index"`
	Device string `json:"device"`
	Serial string `json:"serial"`
	Vendor string `json:"vendor"`
	Total  uint64 `json:"total"`
	Format string `json:"format"`
}

type HardwareInfo struct {
	CpuInfo  []CpuInfo
	NetInfo  []NetInfo
	DiskInfo []DiskInfo
}

func GetCpuInfo() ([]CpuInfo, error) {
	cpus, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	list := []CpuInfo{}
	for _, c := range cpus {
		list = append(list, CpuInfo{
			Index:      c.CPU,
			CoreNum:    c.Cores,
			VemdorId:   c.VendorID,
			PhysicalId: c.PhysicalID,
			ModelName:  c.ModelName,
			MHZ:        c.Mhz,
		})
	}
	return list, nil
}

func GetNetInfo() ([]NetInfo, error) {
	nets, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	list := []NetInfo{}
	for _, inf := range nets {
		if xnet.IsVirtualInterfaceName(inf.Name) {
			continue
		}
		if inf.HardwareAddr == "" {
			continue
		}
		list = append(list, NetInfo{
			Index: int32(inf.Index),
			Name:  inf.Name,
			MTU:   int32(inf.MTU),
			MAC:   inf.HardwareAddr,
		})
	}
	return list, nil
}

func GetDiskInfo() ([]DiskInfo, error) {
	disks, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}
	list := []DiskInfo{}
	for index, part := range disks {
		usage, _ := disk.Usage(part.Mountpoint)
		serial, _ := disk.SerialNumber(part.Device)
		label, _ := disk.Label(part.Device)
		list = append(list, DiskInfo{
			Index:  int32(index),
			Device: part.Device,
			Serial: serial,
			Vendor: label,
			Total:  usage.Total,
			Format: part.Fstype,
		})
	}
	return list, nil
}

func GetHardwareInfo() (*HardwareInfo, error) {
	cpu, _ := GetCpuInfo()
	net, _ := GetNetInfo()
	disk, _ := GetDiskInfo()
	return &HardwareInfo{
		CpuInfo:  cpu,
		NetInfo:  net,
		DiskInfo: disk,
	}, nil
}
