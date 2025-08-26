package discovery

import (
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// calculate weight through cpu/mem
func getDynamicWeight() string {
	// get cpu usage
	percent, err := cpu.Percent(time.Second, false) // 取 1 秒平均 CPU 使用率
	if err != nil || len(percent) == 0 {
		return "1" // minimum score
	}
	cpuUsage := percent[0] // single

	// get mem usage
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "1"
	}
	memUsage := vmStat.UsedPercent

	// cpu/ mem  50%:50% to get score
	score := (100-cpuUsage)*0.5 + (100-memUsage)*0.5

	// range 1-100
	if score < 1 {
		score = 1
	}
	if score > 100 {
		score = 100
	}

	return fmt.Sprintf("%d", int64(math.Round(score)))
}
