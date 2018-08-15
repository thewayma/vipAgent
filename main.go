package main

import (
	"flag"
	"github.com/thewayma/vipAgent/cmdExecute"
	"github.com/thewayma/vipAgent/etcdClient"
	"github.com/thewayma/vipAgent/g"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()

	g.ParseConfig(*cfg)

	cmdExecute.CmdExecute()

	_ = etcdClient.NewWatcher(g.Config().EtcdAddList)

	select {}
}
