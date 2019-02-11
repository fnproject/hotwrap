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
	"strings"
	"syscall"
	"time"

	"github.com/fnproject/fdk-go"
)

var debug = false

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Failed to start hotwrap, no command specified in arguments ")
	}

	if os.Getenv("FN_HOTWRAP_VERBOSE") == "true" {
		debug = true
	}

	cmd := os.Args[1]
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[2:]
	}

	if debug {
		log.Printf("hotwrap running command:  %v %v", cmd, strings.Join(args, " "))
	}
	fdk.Handle(withError(cmd, args))

}

func withError(execName string, execArgs []string) fdk.HandlerFunc {
	f := func(ctx context.Context, in io.Reader, out io.Writer) {
		err := runExec(ctx, fmt.Sprintf(
			"%v %v", execName, strings.Join(execArgs, " ")), in, out)
		if err != nil {
			fdk.WriteStatus(out, http.StatusInternalServerError)
			_, writeErr := io.WriteString(out, err.Error())
			if writeErr != nil && debug {
				log.Print("Failed to write error details to stdout")
			}

			return
		}
		fdk.WriteStatus(out, http.StatusOK)
	}
	return f
}

func runExec(ctx context.Context, execCMDwithArgs string, in io.Reader, out io.Writer) error {
	log.Println(execCMDwithArgs)
	fctx := fdk.GetContext(ctx)
	defer timeTrack(time.Now(), fmt.Sprintf("run-exec-%v", fctx.CallID()))
	cancel := make(chan os.Signal, 3)
	signal.Notify(cancel, os.Interrupt)
	defer signal.Stop(cancel)
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", execCMDwithArgs)
	cmd.Env = os.Environ()
	if in != nil {
		cmd.Stdin = in
	}

	if out != nil {
		cmd.Stdout = out
	}

	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if debug {
					log.Printf("hotwrap: exit code: %d\n", status.ExitStatus())
				}
			}
		}
		return fmt.Errorf("error running exec: %v", err)
	}
	return nil
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	if debug {
		log.Printf("%s took %s\n", name, elapsed)
	}
}
