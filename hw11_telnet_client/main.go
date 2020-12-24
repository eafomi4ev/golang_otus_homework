package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/pflag"
)

var (
	timeout string
)

func init() {
	pflag.StringVar(&timeout, "timeout", "", "timeout")
	pflag.Parse()
}

func main() {
	if pflag.NArg() < 2 {
		log.Fatal("not enough arguments")
	}

	host := pflag.Arg(0)
	port := pflag.Arg(1)
	address := net.JoinHostPort(host, port)

	timeoutD, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatal("Incorrect timeout value")
	}

	telnet := NewTelnetClient(address, timeoutD, os.Stdin, os.Stdout)

	err = telnet.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer telnet.Close()

	ctx, ctxCancelF := context.WithCancel(context.Background())
	go func() {
		defer ctxCancelF()

		err := telnet.Send()
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		defer ctxCancelF()

		err := telnet.Receive()
		if err != nil {
			fmt.Println(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
	case <-ctx.Done():
	}
}
