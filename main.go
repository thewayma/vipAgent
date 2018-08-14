package main

import (
	"flag"
	"github.com/thewayma/vipAgent/etcdClient"
	"github.com/thewayma/vipAgent/g"
	"github.com/thewayma/vipAgent/cmdExecute"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()

	g.ParseConfig(*cfg)

	_ = etcdClient.NewWatcher(g.Config().EtcdAddList)

	cmdExecute.CmdExecute()

	select {}
}
