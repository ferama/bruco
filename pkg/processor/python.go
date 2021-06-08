package processor

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ferama/bruco/pkg/channel"
	"github.com/ferama/bruco/pkg/common"
)

// Python class
type Python struct {
	cmd *exec.Cmd
	ch  *channel.Channel

	name          string
	available     chan *Python
	eventResponse chan Response
}

// NewPython creates a python instance
func NewPython(name string,
	availableWorkers chan *Python,
	wrapperPath string,
	handlerPath string,
	moduleName string,
	env []EnvVar, workingDir string) *Python {

	ch, _ := channel.NewChannel()
	pythonPath, err := common.FindPython()
	if err != nil {
		log.Fatalln("can't find python executable")
	}
	if !filepath.IsAbs(handlerPath) {
		handlerPath = filepath.Join(workingDir, handlerPath)
	}

	args := []string{
		pythonPath, "-u", wrapperPath,
		"--handler-path", handlerPath,
		"--socket", ch.SocketPath,
		"--worker-name", name,
		"--module-name", moduleName,
	}

	cmd := exec.Command(args[0], args[1:]...)

	// builds command environment
	cmd.Env = os.Environ()
	for _, e := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", e.Name, e.Value))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("[PYTHON] %s", err)
	}

	python := &Python{
		cmd:           cmd,
		ch:            ch,
		name:          name,
		available:     availableWorkers,
		eventResponse: make(chan Response),
	}

	go python.handleOutput()

	availableWorkers <- python

	return python
}

func (p *Python) handleOutput() {
	for {
		out, rerr := p.ch.Read()
		if out != nil {
			// the worker ended its job. Mark it as available again
			// putting in the availabe channel
			p.available <- p
		}
		if rerr != nil {
			// log.Printf("[PYTHON] error %s", rerr)
			res := &Response{
				Data:  "",
				Error: rerr.Error(),
			}
			p.eventResponse <- *res
			return
		}

		res := &Response{}
		err := json.Unmarshal(out, res)
		if err != nil {
			log.Printf("[PYTHON] error %s", err)
			res = &Response{
				Data:  "",
				Error: err.Error(),
			}
			p.eventResponse <- *res
			return
		}
		p.eventResponse <- *res
	}
}

// HandleEvent writes to python process socket channel
func (p *Python) handleEvent(data []byte) error {
	return p.ch.Write(data)
}

// Kill kills the python sub process
func (p *Python) kill() {
	if err := p.cmd.Process.Kill(); err != nil {
		log.Fatalln(err)
	}
	// prevents zombie process
	p.cmd.Process.Wait()
	p.ch.Close()
}
