package main

import (
	"bootmanager/cmd"
	_ "net/http/pprof"

	// "github.com/sirupsen/logrus"
)

// func init() {

// 	// logrus.SetFormatter(&logrus.TextFormatter{
// 	// 	FullTimestamp: true,
// 	// })
// 	logrus.SetReportCaller(true)
// }

func main() {
	cmd.Execute()
}