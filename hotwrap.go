package main

import (
	"github.com/fnproject/fdk-go"
	"context"
	"io"
	"os"
	"log"
	"os/exec"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Failed to start hotwrap, no command specified in arguments ")
	}


	if os.Getenv("HOTWRAP_VERBOSE") != "" {

	}
	cmd := os.Args[1]
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[2:]
	}

	fdk.Handle(&hotWrap{
		cmd:  cmd,
		args: args,
		env:  os.Environ(),
	})

}

type hotWrap struct {
	verbose bool
	cmd     string
	args    []string
	env     []string
}

func (hw *hotWrap) logf(fmt string, args ...interface{}) {
	if hw.verbose {
		log.Printf(fmt, args...)
	}
}

func (hw *hotWrap) Serve(ctx context.Context, r io.Reader, w io.Writer) {


	hw.logf("Running '%s %s'", hw.cmd, strings.Join(hw.args," "))

	baseEnv := hw.env

	cmd := exec.Command(hw.cmd, hw.args...)
	cmd.Env = baseEnv
	cmd.Stdout = w
	cmd.Stdin = r

	stderr, err := cmd.StderrPipe()

	if err !=nil {
		log.Fatalf("Failed to open stderr pipe %s",err)

	}

	go func(){
		io.Copy(os.Stderr,stderr)
	}()

	err = cmd.Start()
	if err !=nil {
		log.Fatalf("Failed to start command %s",err)
	}

	err = cmd.Wait()

	if ee,ok:= err.(*exec.ExitError); ok  {
		log.Printf("Command %s exited with status %s",hw.cmd,ee.ProcessState)
		fdk.WriteStatus(w, 500)
	}

}
