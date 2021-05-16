package python

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/ferama/coreai/pkg/channel"
)

// Python class
type Python struct {
	cmd          *exec.Cmd
	killingMutex sync.Mutex
	ch           *channel.Channel

	name      string
	available chan *Python
}

// NewPython creates a python instance
func NewPython(name string, availableWorkers chan *Python,
	wrapperPath string, lambdaPath string) *Python {
	ch, _ := channel.NewChannel()

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
		cmd:       cmd,
		ch:        ch,
		name:      name,
		available: availableWorkers,
	}

	go python.handleOutput()

	availableWorkers <- python

	return python
}

func (p *Python) handleOutput() {
	for {
		out, _ := p.ch.Read()
		if out != nil {
			p.available <- p
		}

		msg := &Message{}
		err := json.Unmarshal(out, msg)
		if err != nil {
			return
		}
		log.Println(fmt.Sprintf("(%s): %s", p.name, msg.Response))
	}
}

// HandleEvent writes to python process stdin
func (p *Python) handleEvent(data []byte) error {
	return p.ch.Write(data)
}

// Kill kills the python sub process
func (p *Python) kill() {
	// this lock is never released
	p.killingMutex.Lock()

	// log.Println("killing python subprocess")
	// Kill it:
	if err := p.cmd.Process.Kill(); err != nil {
		log.Fatalln(err)
	}
	// prevents zombie process
	p.cmd.Process.Wait()
	p.ch.Close()
}
