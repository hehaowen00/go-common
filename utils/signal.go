package utils

import (
	"fmt"
	"os"
	"os/signal"
)

func WaitForInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Print("\r")
}
