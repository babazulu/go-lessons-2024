package main

import (
	"balun-school-task-3/batch"
	"flag"
	"fmt"
	"time"
)

var ip = flag.String("ip", "127.0.0.1", "listen ip")
var port = flag.String("port", "8080", "listen port")
var size = flag.Int("size", 3, "batch size")
var timeout = flag.Int("timeout", 3, "timeout in seconds")

func init() {
	flag.Parse()
}

type Db struct {
}

func (d *Db) SaveBatch(data []string, size int) {
	for _, v := range data {
		fmt.Println(v)
	}
}

func main() {
	ch := make(chan bool)

	b := batch.NewService(&Db{}, *ip, *port, *size, time.Duration(*timeout)*time.Second)
	b.Start()

	fmt.Println("Kill me")
	<-ch
}
