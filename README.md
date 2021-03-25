# nhn_rds_mssql
# build
go build -o nhn_rds_mssql.exe cmd/main.go

# use
telegraf.conf
```
[[inputs.execd]]
  command = ["C:\\Program Files\\Telegraf\\nhn_rds_mssql.exe", "-config", "C:\\Program Files\\Telegraf\\nhn_rds_mssql.conf"]
  signal = "none"
```

nhn_rds_mssql.conf
```
[[inputs.nhn_rds_mssql]]
  taginclude = ["host", "unit", "counter"]
  servers = ["Server=localhost;Port=1433;User Id=rdsadmin;Password=rdsTC20!^;app name=telegraf;log=1;"]
```
