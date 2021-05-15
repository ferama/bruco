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

	name string
	wq   chan *Python
}

// NewPython creates a python instance
func NewPython(name string, wq chan *Python) *Python {

	ch, _ := channel.NewChannel()

	pythonPath := "python3"
	wrapperPath := "./wrapper"
	cmd := exec.Command(pythonPath, wrapperPath, fmt.Sprintf("%d", ch.Port))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	python := &Python{
		cmd:  cmd,
		ch:   ch,
		name: name,
		wq:   wq,
	}

	go python.handleOutput()

	wq <- python

	return python
}

func (p *Python) handleOutput() {
	for {
		out, _ := p.ch.Read()
		if out != nil {
			p.wq <- p
		}

		// log.Print(string(out))
		msg := &Message{}
		err := json.Unmarshal(out, msg)
		if err != nil {
			// log.Println(err)
			return
		}
		log.Println(fmt.Sprintf("(%s): %s", p.name, msg.Response))
	}
}

// HandleEvent writes to python process stdin
func (p *Python) HandleEvent(data []byte) error {
	return p.ch.Write(data)
}

// Kill kills the python sub process
func (p *Python) Kill() {
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
