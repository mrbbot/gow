package main

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

var (
	cmd *exec.Cmd
)

func run(args []string) error {
	if cmd != nil {
		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			syscall.Kill(-pgid, syscall.SIGTERM)
		}
		_ = cmd.Wait()
	}
	cmd = exec.Command("go", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Start()
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Printf("[GOW] not enough args: usage: gow <file>\n")
		os.Exit(1)
	}

	runArgs := append([]string{"run"}, args...)

	err := run(runArgs)
	if err != nil {
		fmt.Printf("[GOW] err starting: %v\n", err)
	}

	w := watcher.New()
	w.FilterOps(watcher.Write, watcher.Create, watcher.Move, watcher.Remove, watcher.Rename)

	go func() {
		for {
			select {
			case event := <-w.Event:
				if filepath.Ext(event.Path) == ".go" {
					fmt.Printf("[GOW] %s changed, restarting...\n", filepath.Base(event.Path))
					err := run(runArgs)
					if err != nil {
						fmt.Printf("[GOW] err restarting: %v\n", err)
					}
				}
			case err := <-w.Error:
				fmt.Printf("[GOW] err watching: %v\n", err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive("."); err != nil {
		fmt.Printf("[GOW] err watching: %v\n", err)
		os.Exit(1)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		fmt.Printf("[GOW] err watching: %v\n", err)
		os.Exit(1)
	}
}
