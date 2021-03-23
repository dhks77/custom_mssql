package custom_mssql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/influxdata/telegraf"
)

func (c *CustomMssql) gatherSqlServer(acc telegraf.Accumulator) error {
	var wg sync.WaitGroup
	query := Query{ScriptName: "SQLServerLogBackupSize-CustomMssql", Script: sqlServerLogBackupSize, ResultByRow: false}

	for _, server := range c.Servers {
		acc.AddError(c.checkServer(server, acc))
		wg.Add(1)
		go func(serv string, query Query) {
			defer wg.Done()
			queryError := c.gatherServer(serv, query, acc)
			acc.AddError(queryError)
		}(server, query)
	}
	wg.Wait()

	return nil
}

func (c *CustomMssql) checkServer(server string, acc telegraf.Accumulator) error {
	var fields = make(map[string]interface{})
	var tags = make(map[string]string)

	db, err := sql.Open("mssql", server)
	if err != nil {
		fields["value"] = 0
		acc.AddFields("sqlserver_connection_is_alive", fields, tags, time.Now())
		return err
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		fields["value"] = 0
		acc.AddFields("sqlserver_connection_is_alive", fields, tags, time.Now())
		return err
	} else {
		fields["value"] = 1
		acc.AddFields("sqlserver_connection_is_alive", fields, tags, time.Now())
	}

	defer db.Close()
	return nil
}

func (c *CustomMssql) gatherServer(server string, query Query, acc telegraf.Accumulator) error {
	// deferred opening
	conn, err := sql.Open("mssql", server)
	if err != nil {
		return err
	}
	defer conn.Close()

	// execute query
	rows, err := conn.Query(query.Script)
	if err != nil {
		return fmt.Errorf("Script %c failed: %w", query.ScriptName, err)
		//return   err
	}
	defer rows.Close()

	// grab the column information from the result
	query.OrderedColumns, err = rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = c.accRow(query, acc, rows)
		if err != nil {
			return err
		}
	}
	return rows.Err()
}

func (c *CustomMssql) accRow(query Query, acc telegraf.Accumulator, row scanner) error {
	var columnVars []interface{}
	var fields = make(map[string]interface{})

	// store the column name with its *interface{}
	columnMap := make(map[string]*interface{})
	for _, column := range query.OrderedColumns {
		columnMap[column] = new(interface{})
	}
	// populate the array of interface{} with the pointers in the right order
	for i := 0; i < len(columnMap); i++ {
		columnVars = append(columnVars, columnMap[query.OrderedColumns[i]])
	}
	// deconstruct array of variables and send to Scan
	err := row.Scan(columnVars...)
	if err != nil {
		return err
	}

	// measurement: identified by the header
	// tags: all other fields of type string
	tags := map[string]string{}
	var measurement string
	for header, val := range columnMap {
		if str, ok := (*val).(string); ok {
			if header == "measurement" {
				measurement = str
			} else {
				tags[header] = str
			}
		}
	}

	if query.ResultByRow {
		// add measurement to Accumulator
		acc.AddFields(measurement,
			map[string]interface{}{"value": *columnMap["value"]},
			tags, time.Now())
	} else {
		// values
		for header, val := range columnMap {
			if _, ok := (*val).(string); !ok {
				fields[header] = *val
			}
		}
		// add fields to Accumulator
		acc.AddFields(pluginName, fields, tags, time.Now())
	}
	return nil
}
