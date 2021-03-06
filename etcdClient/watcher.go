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
	//serviceName map[string]bool //!< /7/idcX/serviceX 标识某业务是否存在
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
		//serviceName: make(map[string]bool),
	}

	watcher.traverseEtcdNodeOnInit()

	go watcher.WatchService()

	return watcher
}

func (m *Watcher) traverseEtcdNodeOnInit() {
	nodeString := fmt.Sprintf("7/%s", g.Config().DefaultTags["Idc"])

	kapi := m.KeysAPI

	resp, err := kapi.Get(context.Background(), nodeString, nil)
	if err != nil {
		log.Fatal(err)
		return
		//!< TODO 容错
	}

	for _, n := range resp.Node.Nodes {
		serviceNodeString := n.Key
		log.Printf("serviceNode=%s\n", serviceNodeString)

		respSub, errSub := kapi.Get(context.Background(), serviceNodeString, nil)
		if errSub != nil {
			log.Printf("get etcd serviceNode=%s, Failure", serviceNodeString)
			continue
		}

		for _, node := range respSub.Node.Nodes {
			keyArray := strings.Split(node.Key, "/")
			keyStr := keyArray[len(keyArray) - 1]

			if keyStr != "vIpPort" {
				continue
			}

			str := strings.Split(node.Value, ":")
			vip := str[0]
			vport := str[1]

			log.Println("traverseEtcdNodeOnInit ETCD set Event: vip =", vip, ", vport =", vport)

			g.AddCh <- vip
		}
	}
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

				log.Println("Watch ETCD set Event: vip =", vip, ", vport =", vport)

				g.AddCh <- vip
			}

		} else if res.Action == "delete" {
			n := res.PrevNode
			keyArray := strings.Split(n.Key, "/")

			if keyArray[len(keyArray)-1] == "vIpPort" {
				str := strings.Split(n.Value, ":")
				vip := str[0]
				vport := str[1]

				log.Println("Watch ETCD delete Event: vip =", vip, ", vport =", vport)

				g.DelCh <- vip
			}

		} else if res.Action == "expire" {
			//!< 目前, 设置节点时不会用到expire属性
		}
	}
}
