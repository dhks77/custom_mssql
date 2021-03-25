package nhn_rds_mssql

import (
	"fmt"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/shirou/gopsutil/cpu"
)

func (c *CustomMssql) gatherCpuUsage(acc telegraf.Accumulator) error {
	times, err := c.ps.CPUTimes(false, true)

	if err != nil {
		return err
	}

	now := time.Now()
	for _, cts := range times {
		tags := map[string]string{
			"cpu": cts.CPU,
		}

		// Add in percentage
		if len(c.lastStats) == 0 {
			// If it's the 1st gather, can't get CPU Usage stats yet
			c.setLastStats(times)
			return nil
		}

		lastCts, ok := c.lastStats[cts.CPU]
		if !ok {
			continue
		}
		total := totalCPUTime(cts)
		lastTotal := totalCPUTime(lastCts)

		totalDelta := total - lastTotal

		if totalDelta < 0 {
			return fmt.Errorf("current total CPU time is less than previous total CPU time")
		}

		if totalDelta == 0 {
			continue
		}

		fields := map[string]interface{}{
			"value": 100 - 100*(cts.Idle-lastCts.Idle)/totalDelta,
		}

		acc.AddGauge("cpu_usage", fields, tags, now)
	}

	c.setLastStats(times)
	return nil
}

func (c *CustomMssql) setLastStats(times []cpu.TimesStat) {
	c.lastStats = make(map[string]cpu.TimesStat)
	for _, cts := range times {
		c.lastStats[cts.CPU] = cts
	}
}

func totalCPUTime(t cpu.TimesStat) float64 {
	total := t.User + t.System + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal + t.Idle
	return total
}
