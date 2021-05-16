package main

import (
	"awesomeProject/src/server"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"runtime"
)

func start() {
	var etcFile = flag.String("c", "", "etc config file")
	var process = flag.String("s", "", "[status]| [shutdown] shutdown the server or check the status")
	flag.Parse()
	bootstrap := server.New()
	if *etcFile == "" {
		if *process == "shutdown" {
			bootstrap.ShutDownServer()
			return
		}
		logrus.Fatal("etc file should not be empty!")
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	bootstrap.Start(*etcFile)
}

func main() {
	start()
}
