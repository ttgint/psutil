package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func handle(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getHardwareData() string {
	runtimeOS := runtime.GOOS
	// memory
	vmStat, err := mem.VirtualMemory()
	handle(err)

	// disk - start from "/" mount point for Linux
	// might have to change for Windows!!
	// don't have a Window to test this out, if detect OS == windows
	// then use "\" instead of "/"

	diskStat, err := disk.Usage("/")
	handle(err)

	percentage, err := cpu.Percent(time.Second, true)
	total, err := cpu.Percent(time.Second, true)
	handle(err)

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	handle(err)

	var r strings.Builder

	r.WriteString("OS : " + runtimeOS + "\n")

	r.WriteString("Total memory: " + strconv.FormatUint(vmStat.Total, 10) + " bytes \n")
	r.WriteString("Free memory: " + strconv.FormatUint(vmStat.Free, 10) + " bytes\n")
	r.WriteString("Percentage used memory: " + strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%\n")

	r.WriteString("Total disk space: " + strconv.FormatUint(diskStat.Total, 10) + " bytes \n")
	r.WriteString("Used disk space: " + strconv.FormatUint(diskStat.Used, 10) + " bytes\n")
	r.WriteString("Free disk space: " + strconv.FormatUint(diskStat.Free, 10) + " bytes\n")
	r.WriteString("Percentage disk space usage: " + strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64) + "%\n")

	// r.WriteString("CPU index number: " + strconv.FormatInt(int64(cpuStat[0].CPU), 10) + "\n")
	// r.WriteString("VendorID: " + cpuStat[0].VendorID + "\n")
	// r.WriteString("Family: " + cpuStat[0].Family + "\n")
	// r.WriteString("Number of cores: " + strconv.FormatInt(int64(cpuStat[0].Cores), 10) + "\n")
	// r.WriteString("Model Name: " + cpuStat[0].ModelName + "\n")
	// r.WriteString("Speed: " + strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64) + " MHz \n")

	r.WriteString("Total CPU usage: " + strconv.FormatFloat(total[0], 'f', 2, 64) + "%\n")
	for idx, cpupercent := range percentage {
		r.WriteString("Current CPU utilization: [" + strconv.Itoa(idx) + "] " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%\n")
	}

	r.WriteString("Hostname: " + hostStat.Hostname + "\n")
	r.WriteString("Uptime: " + strconv.FormatUint(hostStat.Uptime, 10) + "\n")
	r.WriteString("Number of processes running: " + strconv.FormatUint(hostStat.Procs, 10) + "\n")

	return r.String()

}

func main() {
	fmt.Println(getHardwareData())
	ioutil.WriteFile(time.Now().Format("20060102_150405.txt"), []byte(getHardwareData()), 0644)
	for range time.Tick(time.Minute * 10) {
		ioutil.WriteFile(time.Now().Format("20060102_150405.txt"), []byte(getHardwareData()), 0644)
	}
}
