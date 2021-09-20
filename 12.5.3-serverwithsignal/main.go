package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	srv := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%v\n", time.Now())
		}),
	}
	ln, err := net.Listen("tcp4", ":8080")
	if err != nil {
		panic(err)
	}
	log.Println("start")
	go srv.Serve(ln)

	<-signals
	err = srv.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
	log.Println("end")
}
