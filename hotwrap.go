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
	"regexp"
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

// We explicitly omit valid headers and restrict to only those that can be trivially converted to env vars
var validHeaderRegex = regexp.MustCompile("[A-Za-z][A-Za-z0-9-_]*")

func runExec(ctx context.Context, execCMDwithArgs string, in io.Reader, out io.Writer) error {
	log.Println(execCMDwithArgs)
	fctx := fdk.GetContext(ctx)
	defer timeTrack(time.Now(), fmt.Sprintf("run-exec-%v", fctx.CallID()))
	cancel := make(chan os.Signal, 3)
	signal.Notify(cancel, os.Interrupt)
	defer signal.Stop(cancel)
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", execCMDwithArgs)

	callEnv := os.Environ()
	seenHeaders := make(map[string]bool)
	for k, vs := range fctx.Header() {
		if validHeaderRegex.MatchString(k) {
			newHeader := strings.Replace(strings.ToUpper(k), "-", "_", -1)

			var envKey string
			if strings.HasPrefix(newHeader, "FN_") {
				envKey = newHeader
			}

			_, gotEnv := os.LookupEnv(envKey)
			// never overwrite an existing env
			if !gotEnv && !seenHeaders[envKey] {
				seenHeaders[envKey] = true
				callEnv = append(callEnv, fmt.Sprintf("%s=%s", envKey, vs[0]))
			}

		} else {
			if debug {
				log.Printf("saw invalid header key :%s", k	)
			}
		}
	}

	callEnv = append(callEnv, fmt.Sprintf("%s=%s", "FN_CALL_ID", fctx.CallID()))
	callEnv = append(callEnv, fmt.Sprintf("%s=%s", "FN_CONTENT_TYPE", fctx.ContentType()))

	if htcx, ok := fctx.(fdk.HTTPContext); ok {
		callEnv = append(callEnv, fmt.Sprintf("%s=%s", "FN_HTTP_REQUEST_URL",htcx.RequestURL()))
		callEnv = append(callEnv, fmt.Sprintf("%s=%s", "FN_HTTP_METHOD",htcx.RequestMethod()))

	}

	cmd.Env = callEnv
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
