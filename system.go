package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/disk"

	"github.com/shirou/gopsutil/mem"

	"github.com/shirou/gopsutil/cpu"
)

const readDuration = time.Second
const tempClassPath = "/sys/class/thermal/thermal_zone0/temp"

type CPUInfo struct {
	Usage []float64 `json:"usage"`
	Temp  int       `json:"temp"`
}

type MemoryInfo struct {
	MemoryUsedPercent float64 `json:"memoryUsedPercent"`
	SwapUsedPercent   float64 `json:"swapUsedPercent"`
}

type DiskInfo struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
}

type SystemInfo struct {
	CPU    *CPUInfo    `json:"cpu"`
	Memory *MemoryInfo `json:"memory"`
	Disk   *DiskInfo   `json:"disk"`
}

var errorLog *log.Logger

func init() {
	errorLog = log.New(os.Stderr, "system", log.LstdFlags)
}

func readCPUInfo() (*CPUInfo, error) {
	pers, err := cpu.Percent(readDuration, true)
	if err != nil {
		return nil, err
	}
	return &CPUInfo{
		pers,
		readTemp(),
	}, nil
}

func readMemoryInfo() (*MemoryInfo, error) {
	m, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	s, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}
	return &MemoryInfo{
		MemoryUsedPercent: m.UsedPercent,
		SwapUsedPercent:   s.UsedPercent,
	}, nil
}

func readDiskInfo() (*DiskInfo, error) {
	usg, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}
	return &DiskInfo{
		Total: usg.Total,
		Used:  usg.Used,
	}, nil
}

func readSystemInfoAll() SystemInfo {
	cpuInfoC := make(chan *CPUInfo)
	memoryInfoC := make(chan *MemoryInfo)
	diskInfoC := make(chan *DiskInfo)
	go func(c chan *CPUInfo) {
		defer close(c)
		cpu, err := readCPUInfo()
		if err != nil {
			errorLog.Println("read cpu info error:", err)
		} else {
			c <- cpu
		}
	}(cpuInfoC)
	go func(c chan *MemoryInfo) {
		defer close(c)
		mem, err := readMemoryInfo()
		if err != nil {
			errorLog.Println("read memory info error:", err)
		} else {
			c <- mem
		}
	}(memoryInfoC)
	go func(c chan *DiskInfo) {
		defer close(c)
		dk, err := readDiskInfo()
		if err != nil {
			errorLog.Println("read disk info error:", err)
		} else {
			c <- dk
		}
	}(diskInfoC)
	return SystemInfo{
		CPU:    <-cpuInfoC,
		Memory: <-memoryInfoC,
		Disk:   <-diskInfoC,
	}
}

func readTemp() int {
	b, err := ioutil.ReadFile(tempClassPath)
	if err != nil {
		return -1
	}
	s := strings.TrimRight(string(b), "\n")
	n, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return n
}
