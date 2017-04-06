package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mkdirTool/app"
	"mkdirTool/public"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")

	flag.Parse()

	public.Init()

	public.ParseConfig(*cfg) //解析配置文件

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	public.WaitGroup.Wrap(func() {
		_ = app.WalkDir(public.Config.Root, 0)
	})

	fmt.Println("File processing...")
	<-exitChan
}
