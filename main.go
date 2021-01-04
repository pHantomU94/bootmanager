package main

import (
	"bootmanager/cmd"
	"fmt"

	"net/http"
	_ "net/http/pprof"

	"github.com/sirupsen/logrus"
)

func init() {

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetReportCaller(true)
}

func main() {
	logrus.Info("start")
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	cmd.Execute()
}