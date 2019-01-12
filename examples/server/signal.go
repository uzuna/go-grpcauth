package main

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitSignal is wait exit like signal
func WaitSignal() os.Signal {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt,
	)
	return <-sigc
}
