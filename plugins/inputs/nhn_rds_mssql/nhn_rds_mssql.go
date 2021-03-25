package nhn_rds_mssql

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/inputs/system"
	"github.com/shirou/gopsutil/cpu"
	"log"
)

const pluginName = "nhn_rds_mssql"

// CustomMssql struct
type CustomMssql struct {
	Servers   []string `toml:"servers"`
	ps        system.PS
	lastStats map[string]cpu.TimesStat
}

// Query struct
type Query struct {
	ScriptName     string
	Script         string
	ResultByRow    bool
	OrderedColumns []string
}

const defaultServer = "Server=.;app name=telegraf;log=1;"

const sampleConfig = `
servers = [
  "Server=192.168.1.10;Port=1433;User Id=<user>;Password=<pw>;app name=telegraf;log=1;",
]
`

// SampleConfig return the sample configuration
func (c *CustomMssql) SampleConfig() string {
	return sampleConfig
}

// Description return plugin description
func (c *CustomMssql) Description() string {
	return "Read metrics for custom mssql metric"
}

type scanner interface {
	Scan(dest ...interface{}) error
}

// Gather collect data from SQL Server
func (c *CustomMssql) Gather(acc telegraf.Accumulator) error {
	err := c.gatherCpuUsage(acc)
	if err != nil {
		return err
	}

	err = c.gatherMem(acc)
	if err != nil {
		return err
	}

	err = c.gatherSqlServer(acc)
	if err != nil {
		return err
	}

	return nil
}

func (c *CustomMssql) Init() error {
	if len(c.Servers) == 0 {
		log.Println("W! Warning: Server list is empty.")
	}

	return nil
}

func init() {
	inputs.Add(pluginName, func() telegraf.Input {
		return &CustomMssql{
			Servers: []string{defaultServer},
			ps:      system.NewSystemPS(),
		}
	})
}
