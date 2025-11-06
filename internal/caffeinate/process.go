package caffeinate

import (
	"os/exec"
	"strconv"
	"log"
	"sync"
	"syscall"
)

var (
	cmd       *exec.Cmd
	cmdLock   sync.Mutex
	isRunning bool
	CaffeinateOptions = []string{"-dims"}
)

func IsRunning() bool {
	cmdLock.Lock()
	defer cmdLock.Unlock()
	return isRunning
}

func Start(options []string) error {
	cmdLock.Lock()
	defer cmdLock.Unlock()

	if isRunning {
		log.Println("Caffeinate already running")
		return nil
	}

	args := []string{}
	for _, opt := range options {
		if opt == "t" {
			// 't' means timeout; handled with next element
			continue
		}
		if _, err := strconv.Atoi(opt); err == nil {
			continue
		}
		// If the option already starts with '-', keep it as-is
		if len(opt) > 0 && opt[0] == '-' {
			args = append(args, opt)
		} else {
			// Add '-' prefix for single-letter options (d, i, m, s, u)
			args = append(args, "-"+opt)
		}
	}
	for i := 0; i < len(options); i++ {
		if options[i] == "t" && i+1 < len(options) {
			args = append(args, "-t", options[i+1])
		}
	}
	log.Printf("Starting caffeinate with args: %v\n", args)
	cmd = exec.Command("caffeinate", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return err
	}

	isRunning = true
	log.Printf("Started caffeinate (PID %d)\n", cmd.Process.Pid)
	go func() {
		_ = cmd.Wait()
		cmdLock.Lock()
		defer cmdLock.Unlock()
		isRunning = false
		log.Printf("caffeinate timedout (PID %d)\n", cmd.Process.Pid)
		cmd = nil
	}()
	return nil
}

func StartTimed(seconds int, options []string) error {
	strSeconds := strconv.Itoa(seconds)
	opts := append(options, "t", strSeconds)
	return Start(opts)
}

func Stop() error {
	cmdLock.Lock()
	defer cmdLock.Unlock()

	if cmd == nil || !isRunning {
		log.Println("No caffeinate process to stop")
		return nil
	}

	if err := cmd.Process.Kill(); err != nil {
		return err
	}

	log.Printf("Stopped caffeinate (PID %d)\n", cmd.Process.Pid)
	isRunning = false
	cmd = nil
	return nil
}

func concatOptions(opts []string) string {
	out := ""
	for _, o := range opts {
		out += o
	}
	return out
}
