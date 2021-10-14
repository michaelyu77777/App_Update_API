package main

import (
	_ "github.com/denisenkom/go-mssqldb"

	"leapsy.com/servers"
)

// main - 主程式
func main() {
	servers.StartServer() // 啟動伺服器
}
