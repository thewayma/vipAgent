package main

import (
	"flag"
	"github.com/coreos/etcd/client"
	"github.com/thewayma/vipAgent/g"
	"time"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()

	g.ParseConfig(*cfg)

	etcdCfg := client.Config{
		Endpoints:               g.Config().EtcdAddList,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(etcdCfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	select {}
}
