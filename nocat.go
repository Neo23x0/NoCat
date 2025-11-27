package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

func main() {
	// nc-like flags
	listen := flag.Bool("l", false, "listen mode (like nc -l)")
	port := flag.Int("p", 0, "port number")
	verbose := flag.Bool("v", false, "verbose output")

	// dummy flags, for nc-style appearance
	keepOpen := flag.Bool("k", false, "dummy keep-open flag (simulated only)")
	execCmd := flag.String("e", "", "dummy exec flag")
	noDNS := flag.Bool("n", false, "dummy no-dns flag (simulated only)")

	flag.Parse()

	if *port == 0 {
		fmt.Fprintln(os.Stderr, "port (-p) is required")
		os.Exit(1)
	}

	// If -e was provided with any value, spawn a fixed, harmless app per OS.
	if *execCmd != "" {
		startDummyChild(*verbose, *execCmd)
	}

	addr := fmt.Sprintf("0.0.0.0:%d", *port)

	if *listen {
		runListenMode(addr, *verbose, *keepOpen, *execCmd, *noDNS)
	} else {
		runClientMode(*port, *verbose)
	}
}

// startDummyChild never executes the value of -e
func startDummyChild(verbose bool, execValue string) {
	cmd := exec.Command(execValue) // placeholder, will be replaced

	// If you want the child to stay alive, attach stdio to the parent:
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Cross-platform: resolves to “new console” on Windows, nil elsewhere
	cmd.SysProcAttr = newConsoleAttr

	if err := cmd.Start(); err != nil {
		if verbose {
			log.Printf("NoCat: failed to start %s: %v", execValue, err)
		}
		return
	}
	if verbose {
		log.Printf("NoCat: %s started with PID %d", execValue, cmd.Process.Pid)
	}
}

func runListenMode(addr string, verbose, keepOpen bool, execCmd string, noDNS bool) {
	if verbose {
		log.Printf(
			"listening on %s (flags: -l -p %d -v -k=%v -e=%q -n=%v)",
			addr, extractPort(addr), keepOpen, execCmd, noDNS,
		)
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen error on %s: %v", addr, err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if verbose {
				log.Printf("accept error: %v", err)
			}
			continue
		}

		if verbose {
			log.Printf("incoming connection from %s", conn.RemoteAddr())
		}

		// Simulate a listener that holds the connection briefly.
		go func(c net.Conn) {
			defer c.Close()
			if verbose {
				log.Printf("holding connection to %s for a short time", c.RemoteAddr())
				if execCmd != "" {
					log.Printf("Note: -e %q was specified; NoCat never executes this program, it only may start a fixed, harmless app depending on the OS.", execCmd)
				}
			}
			time.Sleep(30 * time.Second)
		}(conn)

		// without -k: single connection then exit
		if !keepOpen {
			if verbose {
				log.Printf("single-connection mode (no -k) - exiting listen loop")
			}
			break
		}
	}

	if verbose {
		log.Printf("listener on %s shutting down", addr)
	}
}

func runClientMode(port int, verbose bool) {
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "target host is required when not using -l")
		os.Exit(1)
	}

	host := flag.Arg(0)
	target := fmt.Sprintf("%s:%d", host, port)

	if verbose {
		log.Printf("connecting to %s", target)
	}

	conn, err := net.Dial("tcp", target)
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer conn.Close()

	if verbose {
		log.Printf("connected to %s, idling briefly then closing", conn.RemoteAddr())
	}

	time.Sleep(10 * time.Second)
}

// extractPort is just for nicer logs.
func extractPort(addr string) int {
	var host string
	var port int
	fmt.Sscanf(addr, "%s:%d", &host, &port)
	return port
}
