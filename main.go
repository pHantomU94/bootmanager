/*
 * @Title:
 * @Description:
 * @Version:
 * @Company: Casia
 * @Author: hsj
 * @Date: 2020-12-15 19:47:14
 * @LastEditors: hsj
 * @LastEditTime: 2021-07-22 11:39:32
 */
package main

import (
	"bootmanager/cmd"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	// "github.com/sirupsen/logrus"
)

func init() {

	// logrus.SetFormatter(&logrus.TextFormatter{
	// 	FullTimestamp: true,
	// })
	// logrus.SetReportCaller(true)
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		select {
		case sig := <-c:
			fmt.Printf("Got %s signal.aborting...\n", sig)
			pid := -os.Getpid()
			fmt.Printf("kill %d\n", pid)
			syscall.Kill(pid, syscall.SIGKILL)

		}
	}()
	cmd.Execute()
}
