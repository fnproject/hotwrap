package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/fnproject/fdk-go"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Failed to start hotwrap, no command specified in arguments ")
	}

	cmd := os.Args[1]
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[2:]
	}

	fdk.Handle(withError(cmd, args))

}

func withError(execName string, execArgs []string) fdk.HandlerFunc {
	f :=  func(ctx context.Context, in io.Reader, out io.Writer){
		err := runExec(ctx, execName, execArgs, in, out)
		if err != nil {
			fdk.WriteStatus(out, http.StatusInternalServerError)
			io.WriteString(out, err.Error())
			return
		}
		fdk.WriteStatus(out, http.StatusOK)
	}
	return f
}

func runExec(ctx context.Context, execName string, execArgs []string, in io.Reader, out io.Writer) error {
	fctx := fdk.GetContext(ctx)
	defer timeTrack(time.Now(), fmt.Sprintf("run-exec-%v", fctx.CallID()))
	cancel := make(chan os.Signal, 3)
	signal.Notify(cancel, os.Interrupt)
	defer signal.Stop(cancel)
	result := make(chan error, 1)
	quit := make(chan struct{})
	cmd := exec.CommandContext(ctx, execName, execArgs...)
	if in != nil {
		cmd.Stdin = in
	}
	if out != nil {
		cmd.Stdout = out
	}
	cmd.Stderr = os.Stderr

	go func(cmd *exec.Cmd, done chan<- error) {
		done <- cmd.Run()
	}(cmd, result)

	select {
	case err := <-result:
		close(quit)
		fmt.Fprintln(os.Stderr)
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					log.Printf("exit code: %d\n", status.ExitStatus())
				}
			}
			return fmt.Errorf("error running exec: %v", err)
		}
	}
	return nil
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s\n", name, elapsed)
}
