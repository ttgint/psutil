package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
}

func GetHardwareData() string {
	runtimeOS := runtime.GOOS
	// memory
	vmStat, err := mem.VirtualMemory()
	dealwithErr(err)

	// disk - start from "/" mount point for Linux
	// might have to change for Windows!!
	// don't have a Window to test this out, if detect OS == windows
	// then use "\" instead of "/"

	diskStat, err := disk.Usage("/")
	dealwithErr(err)

	// cpu - get CPU number of cores and speed
	cpuStat, err := cpu.Info()
	dealwithErr(err)
	percentage, err := cpu.Percent(0, true)
	dealwithErr(err)

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	dealwithErr(err)

	// get interfaces MAC/hardware address
	dealwithErr(err)

	var r strings.Builder

	r.WriteString("OS : " + runtimeOS + "\n")

	r.WriteString("Total memory: " + strconv.FormatUint(vmStat.Total, 10) + " bytes \n")
	r.WriteString("Free memory: " + strconv.FormatUint(vmStat.Free, 10) + " bytes\n")
	r.WriteString("Percentage used memory: " + strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%\n")

	// get disk serial number.... strange... not available from disk package at compile time
	// undefined: disk.GetDiskSerialNumber
	//serial := disk.GetDiskSerialNumber("/dev/sda")

	//r.WriteString("Disk serial number: " + serial + "\n"

	r.WriteString("Total disk space: " + strconv.FormatUint(diskStat.Total, 10) + " bytes \n")
	r.WriteString("Used disk space: " + strconv.FormatUint(diskStat.Used, 10) + " bytes\n")
	r.WriteString("Free disk space: " + strconv.FormatUint(diskStat.Free, 10) + " bytes\n")
	r.WriteString("Percentage disk space usage: " + strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64) + "%\n")

	// since my machine has one CPU, I'll use the 0 index
	// if your machine has more than 1 CPU, use the correct index
	// to get the proper data
	r.WriteString("CPU index number: " + strconv.FormatInt(int64(cpuStat[0].CPU), 10) + "\n")
	r.WriteString("VendorID: " + cpuStat[0].VendorID + "\n")
	r.WriteString("Family: " + cpuStat[0].Family + "\n")
	r.WriteString("Number of cores: " + strconv.FormatInt(int64(cpuStat[0].Cores), 10) + "\n")
	r.WriteString("Model Name: " + cpuStat[0].ModelName + "\n")
	r.WriteString("Speed: " + strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64) + " MHz \n")

	for idx, cpupercent := range percentage {
		r.WriteString("Current CPU utilization: [" + strconv.Itoa(idx) + "] " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%\n")
	}

	r.WriteString("Hostname: " + hostStat.Hostname + "\n")
	r.WriteString("Uptime: " + strconv.FormatUint(hostStat.Uptime, 10) + "\n")
	r.WriteString("Number of processes running: " + strconv.FormatUint(hostStat.Procs, 10) + "\n")

	// another way to get the operating system name
	// both darwin for Mac OSX, For Linux, can be ubuntu as platform
	// and linux for OS

	r.WriteString("OS: " + hostStat.OS + "\n")
	r.WriteString("Platform: " + hostStat.Platform + "\n")

	// the unique hardware id for this machine
	r.WriteString("Host ID(uuid): " + hostStat.HostID + "\n")

	return r.String()

}

func main() {
	for range time.Tick(time.Minute * 10) {
		ioutil.WriteFile(time.Now().Format("20060102_150405"), []byte(GetHardwareData()), 0644)
	}
}
