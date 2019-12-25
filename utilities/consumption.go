package utilities

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/supudo/Kuplung-Go/settings"
)

// GetUsage will return overall stats
func GetUsage() string {
	usage := ""
	mem := GetMemoryUsage(true)
	cpu := GetCPUUsage(true)
	if len(mem) == 0 {
		usage += "Memory: n/a"
	} else {
		usage += mem
	}
	usage += ", "
	if len(cpu) == 0 {
		usage += "CPU: n/a"
	} else {
		usage += cpu
	}
	return usage
}

// GetMemoryUsage will return formatted string of application memory allocaitons
func GetMemoryUsage(formatted bool) string {
	sett := settings.GetSettings()
	sett.Consumption.ConsumptionTimerMemory = getMilliseconds()
	if sett.Consumption.ConsumptionIntervalMemory > 0 && sett.Consumption.ConsumptionTimerMemory > sett.Consumption.ConsumptionCounterMemory+(1000*sett.Consumption.ConsumptionIntervalMemory) {
		sett.Consumption.ConsumptionMemory = getMemUsage(formatted)
		sett.Consumption.ConsumptionCounterMemory = sett.Consumption.ConsumptionTimerMemory
	}
	return sett.Consumption.ConsumptionMemory
}

func getMemUsage(formatted bool) string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if formatted {
		return fmt.Sprintf("Memory: %v MiB", bToMb(m.Alloc))
	}
	return fmt.Sprintf("%v", m.Alloc)
}

// GetCPUUsage will return formatted string of application CPU load
func GetCPUUsage(formatted bool) string {
	sett := settings.GetSettings()
	sett.Consumption.ConsumptionTimerCPU = getMilliseconds()
	if sett.Consumption.ConsumptionIntervalCPU > 0 && sett.Consumption.ConsumptionTimerCPU > sett.Consumption.ConsumptionCounterCPU+(1000*sett.Consumption.ConsumptionIntervalCPU) {
		sett.Consumption.ConsumptionCPU = getCPU(formatted)
		sett.Consumption.ConsumptionCounterCPU = sett.Consumption.ConsumptionTimerCPU
	}
	return sett.Consumption.ConsumptionCPU
}

func getCPU(formatted bool) string {
	pid := os.Getpid()
	cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", pid), "-o", "pcpu")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		settings.LogWarn("[Consumption] Can't get process information: ", err)
	}
	_, _ = out.ReadString('\n')
	line, err := out.ReadString('\n')
	if err != nil {
		settings.LogWarn("[Consumption] Can't read process information: ", err)
	}
	if formatted {
		return fmt.Sprintf("CPU: %v%%", strings.TrimSpace(line))
	}
	return fmt.Sprintf("%v", strings.TrimSpace(line))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func getMilliseconds() int64 {
	now := time.Now()
	unixNano := now.UnixNano()
	return unixNano / 1000000
}
