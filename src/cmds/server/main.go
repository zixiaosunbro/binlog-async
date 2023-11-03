package main

import (
	"binlog-async/src"
	"binlog-async/src/application"
)

func main() {
	// redis mysql connection initialize
	src.ResourceInit()
	// async service initialize
	go application.InitAsyncSvc()
	// mysql binlog initialize
	application.InitBinlogSvc()
}
