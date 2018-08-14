package cmdExecute

import (
	"log"
	"github.com/thewayma/vipAgent/g"
)

func CmdExecute() {
	go func() {
		for i := range g.AddCh {
			log.Println("Add vip=", i)
		}
	}()


	go func() {
		for i := range g.DelCh {
			log.Println("Del vip=", i)
		}

	}()
}