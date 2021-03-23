# custom_mssql
# build
go build -o custom_mssql.exe cmd/main.go

# use
telegraf.conf
```
[[inputs.execd]]
  command = ["C:\\Program Files\\Telegraf\\custom_mssql.exe", "-config", "C:\\Program Files\\Telegraf\\custom_mssql.conf"]
  signal = "none"
```

custom_mssql.conf
```
[[inputs.custom_mssql]]
  taginclude = ["host", "unit", "counter"]
  servers = ["Server=localhost;Port=1433;User Id=rdsadmin;Password=rdsTC20!^;app name=telegraf;log=1;"]
```
