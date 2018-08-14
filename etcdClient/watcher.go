package etcdClient

import (
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/thewayma/vipAgent/g"
	"golang.org/x/net/context"
	"log"
	"sort"
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

			/**
			// fetch directory
			resp, err = kapi.Get(context.Background(), nodeString, nil)
			if err != nil {
				log.Fatal(err)
			}*/

			// print directory keys
			sort.Sort(res.Node.Nodes)
			for _, n := range res.Node.Nodes {
				fmt.Printf("Key: %q, Value: %q\n", n.Key, n.Value)
			}

		} else if res.Action == "delete" {

		} else if res.Action == "expire" {
			//!< 目前, 设置节点时不会用到expire属性
		}
	}
}
