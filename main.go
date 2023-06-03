package main

import (
	"flag"

	"github.com/arisnotargon/squad_a2s_cache/funcs"
	"github.com/davecgh/go-spew/spew"
)

var (
	runMode string
)

func init() {
	spew.Dump("in init")
	flag.StringVar(&runMode, "runMode", "listen", "运行模式,list_dev 列出设备名,listen 开始监听抓包")
}

func main() {
	flag.Parse()

	switch runMode {
	case "list_dev":
		funcs.ListDev()
	case "listen":
		funcs.Cap_a2s()
	case "send_test_payload":
		funcs.SendTestPackage()
	case "start_test_server_info":
		funcs.StartTestServerInfo()
	}

}
