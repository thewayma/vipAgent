package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"sync"
)

var (
	config *GlobalConfig
	lock   = new(sync.RWMutex)
)

type GlobalConfig struct {
	NicName     []string
	EtcdAddList []string
	DefaultTags map[string]string
}

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Println("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Println("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Println("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Println("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	config = &c
	lock.Unlock()

	log.Println("read config file:", cfg, "successfully")
}
