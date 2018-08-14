package g

import (
	"runtime"
)

var AddCh = make(chan string)
var DelCh = make(chan string)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

}
