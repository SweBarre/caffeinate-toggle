package caffeinate

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var (
	mu        sync.Mutex
	cmd       *exec.Cmd
	isRunning bool
)

// Start starts the caffeinate process with provided arguments.
func Start(args string) error {
	mu.Lock()
	defer mu.Unlock()

	if cmd != nil && cmd.Process != nil {
		log.Println("Start called, but caffeinate is already running")
		return nil // Already running
	}

	argList := strings.Fields(args)
	cmd = exec.Command("caffeinate", argList...)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start caffeinate: %v\n", err)
	}
	log.Printf("Starting caffeinate (PID %d) with args: %v\n", cmd.Process.Pid, argList)
	isRunning = true
	go func() {
		_ = cmd.Wait()
		mu.Lock()
		defer mu.Unlock()
		isRunning = false
		log.Printf("caffeinate timedout (PID %d)\n", cmd.Process.Pid)
		cmd = nil
	}()
	return nil
}

// StartTimed starts caffeinate with provided args and a timeout in seconds.
func StartTimed(args string, seconds int) error {
	argsWithTimeout := args + " -t " + strconv.Itoa(seconds)
	return Start(argsWithTimeout)
}

// Stop stops the caffeinate process if it's running.
func Stop() error {
	mu.Lock()
	defer mu.Unlock()

	if cmd == nil || isRunning == false {
		return nil // Not running
	}

	err := cmd.Process.Kill()
	if err != nil {
		log.Printf("Failed to kill caffeinate (PID %d)\n", cmd.Process.Pid)
		return err
	}
	log.Printf("Stopped caffeinate (PID %d)\n", cmd.Process.Pid)
	isRunning = false
	return err
}

// IsRunning checks if the caffeinate process is running.
func IsRunning() bool {
	mu.Lock()
	defer mu.Unlock()
	return isRunning
}
