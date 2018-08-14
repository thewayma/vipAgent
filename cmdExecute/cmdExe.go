package cmdExecute

import (
	"log"
	"github.com/thewayma/vipAgent/g"
	"fmt"
	"os/exec"
)

func execute(cmd string, comment string) {
	_, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		log.Printf("%s %s %s\n", comment, cmd, "Failure")
	} else {
		log.Printf("%s %s %s\n", comment, cmd, "Success")
	}
}

func CmdExecute() {
	go func() {
		for i := range g.AddCh {
			log.Println("AddChannel received addVip =", i)

			addIp := fmt.Sprintf("ip address add %s/32 dev lo", i)
			addRoute := fmt.Sprintf("route add -host %s dev lo", i)

			execute(addIp, "Add vip:")
			execute(addRoute, "Add route:")
		}
	}()

	go func() {
		for i := range g.DelCh {
			log.Println("DelChannel received delVip =", i)

			delIp := fmt.Sprintf("ip address del %s/32 dev lo", i)
			delRoute := fmt.Sprintf("route del -host %s dev lo", i)

			execute(delIp, "Del vip:")
			execute(delRoute, "Del route:")
		}
	}()
}