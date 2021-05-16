package pool

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
	eventResponse chan message
}

// NewPython creates a python instance
func NewPython(name string, availableWorkers chan *Python,
	wrapperPath string, lambdaPath string) *Python {
	ch, _ := newChannel()

	pythonPath := "python3"

	args := []string{
		pythonPath, "-u", wrapperPath,
		"--lambda-path", lambdaPath,
		"--port", fmt.Sprintf("%d", ch.Port),
		"--worker-name", name,
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	python := &Python{
		cmd:           cmd,
		ch:            ch,
		name:          name,
		available:     availableWorkers,
		eventResponse: make(chan message),
	}

	go python.handleOutput()

	availableWorkers <- python

	return python
}

func (p *Python) handleOutput() {
	for {
		out, rerr := p.ch.Read()
		if out != nil {
			p.available <- p
		}
		if rerr != nil {
			return
		}

		msg := &message{}
		err := json.Unmarshal(out, msg)
		if err != nil {
			return
		}
		// p.eventResponse <- *msg
		log.Println(fmt.Sprintf("(%s): %s", p.name, msg.Response))
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