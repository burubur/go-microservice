package servd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/kudoindonesia/frontier_biller_sdk/log"
)

func shutdown(processChannel chan bool, interuptionSignal chan os.Signal) {
	signal.Notify(interuptionSignal, syscall.SIGINT, syscall.SIGTERM)

	terminationSignal := <-interuptionSignal
	log.Warn(fmt.Sprint("Caught signal: ", terminationSignal))
	signal.Stop(interuptionSignal)
	processChannel <- true
}
