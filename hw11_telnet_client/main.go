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
	timeout time.Duration
)

func init() {
	pflag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout")
	pflag.Parse()
}

func main() {
	if pflag.NArg() < 2 {
		log.Fatal("not enough arguments")
	}

	host := pflag.Arg(0)
	port := pflag.Arg(1)
	address := net.JoinHostPort(host, port)

	telnet := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := telnet.Connect(); err != nil {
		log.Fatal(err)
	}
	defer telnet.Close()

	ctx, ctxCancelF := context.WithCancel(context.Background())
	go func() {
		defer ctxCancelF()

		err := telnet.Send()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	go func() {
		defer ctxCancelF()

		err := telnet.Receive()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
	case <-ctx.Done():
		close(sigCh)
	}
}
