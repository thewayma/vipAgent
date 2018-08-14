package etcdClient

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/thewayma/vipAgent/g"
	"log"
	"strings"
	"time"
)

type Watcher struct {
	KeysAPI     client.KeysAPI  //!< etcd client
	serviceName map[string]bool //!< /7/idcX/serviceX 标识某业务是否存在
}

func NewWatcher(endpoints []string) *Watcher {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	watcher := &Watcher{
		KeysAPI:     client.NewKeysAPI(etcdClient),
		serviceName: make(map[string]bool),
	}

	go watcher.WatchService()

	return watcher
}

func (m *Watcher) WatchService() {
	nodeString := fmt.Sprintf("7/%s", g.Config().DefaultTags["Idc"])

	kapi := m.KeysAPI
	watcher := kapi.Watcher(nodeString, &client.WatcherOptions{
		Recursive: true,
	})

	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch workers:", err)
			break
		}

		if res.Action == "set" {
			n := res.Node
			keyArray := strings.Split(n.Key, "/")

			if keyArray[len(keyArray)-1] == "vIpPort" {
				str := strings.Split(n.Value, ":")
				vip := str[0]
				vport := str[1]

				log.Println("Watch ETCD set Event: vip=", vip, ", vport=", vport)

				g.AddCh <- vip
			}

		} else if res.Action == "delete" {
			n := res.Node
			keyArray := strings.Split(n.Key, "/")

			if keyArray[len(keyArray)-1] == "vIpPort" {
				str := strings.Split(n.Value, ":")
				vip := str[0]
				vport := str[1]

				log.Println("Watch ETCD delete Event: vip=", vip, ", vport=", vport)

				g.DelCh <- vip
			}

		} else if res.Action == "expire" {
			//!< 目前, 设置节点时不会用到expire属性
		}
	}
}
