package processor

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Python class
type Python struct {
	cmd *exec.Cmd
	ch  *channel

	name          string
	available     chan *Python
	eventResponse chan Response
}

// NewPython creates a python instance
func NewPython(name string, availableWorkers chan *Python,
	wrapperPath string, workdir string, moduleName string) *Python {
	ch, _ := newChannel()

	pythonPath := "python3"

	args := []string{
		pythonPath, "-u", wrapperPath,
		"--workdir", workdir,
		"--port", fmt.Sprintf("%d", ch.Port),
		"--worker-name", name,
		"--module-name", moduleName,
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
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
			return
		}

		res := &Response{}
		err := json.Unmarshal(out, res)
		if err != nil {
			return
		}
		p.eventResponse <- *res
	}
}

// HandleEvent writes to python process stdin
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
