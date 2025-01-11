package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [--timeout=<duration>] <host> <port>")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := fmt.Sprintf("%s:%s", host, port)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Fprintln(os.Stderr, "...Terminating")
		client.Close()
		os.Exit(0)
	}()

	done := make(chan error, 2)
	go func() { done <- client.Send() }()
	go func() { done <- client.Receive() }()

	if err := <-done; err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	fmt.Fprintln(os.Stderr, "...Connection closed")
}
