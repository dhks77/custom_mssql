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