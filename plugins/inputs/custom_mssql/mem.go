package custom_mssql

import (
	"fmt"
	"github.com/influxdata/telegraf"
)

func (c *CustomMssql) gatherMem(acc telegraf.Accumulator) error {
	vm, err := c.ps.VMStat()
	if err != nil {
		return fmt.Errorf("error getting virtual memory info: %s", err)
	}

	fields := map[string]interface{}{
		"value": 100 - 100 * float64(vm.Available) / float64(vm.Total),
	}

	acc.AddGauge("memory_usage", fields, nil)
	return nil
}